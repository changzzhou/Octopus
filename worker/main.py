"""
Octopus Worker - Python Task Executor

This is the main entry point for the Python worker. It connects to Redis,
listens for task requests, and executes them using registered handlers.
"""

import asyncio
import logging
import os
import signal
import sys
from typing import NoReturn

from dotenv import load_dotenv

from protocol import TaskRequest, TaskResponse, TaskStatus
from handlers import get_handler, list_handlers

load_dotenv()

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger("octopus-worker")

REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379")
TASK_QUEUE = os.getenv("TASK_QUEUE", "octopus:tasks")
RESULT_QUEUE = os.getenv("RESULT_QUEUE", "octopus:results")
WORKER_ID = os.getenv("WORKER_ID", f"worker-{os.getpid()}")


class Worker:
    """Octopus task worker."""
    
    def __init__(self, redis_url: str, task_queue: str, result_queue: str):
        self.redis_url = redis_url
        self.task_queue = task_queue
        self.result_queue = result_queue
        self.running = False
        self._redis = None
    
    async def connect(self) -> None:
        """Connect to Redis."""
        try:
            import redis.asyncio as redis
            self._redis = redis.from_url(self.redis_url)
            await self._redis.ping()
            logger.info(f"Connected to Redis: {self.redis_url}")
        except ImportError:
            logger.warning("Redis not available - running in standalone mode")
            self._redis = None
        except Exception as e:
            logger.warning(f"Failed to connect to Redis: {e} - running in standalone mode")
            self._redis = None
    
    async def disconnect(self) -> None:
        """Disconnect from Redis."""
        if self._redis:
            await self._redis.close()
            logger.info("Disconnected from Redis")
    
    async def process_task(self, request: TaskRequest) -> TaskResponse:
        """Process a single task."""
        handler = get_handler(request.task_type)
        if handler is None:
            logger.error(f"No handler for task type: {request.task_type}")
            return TaskResponse.failure(
                request.task_id,
                f"Unknown task type: {request.task_type}",
                request.context,
            )
        
        try:
            response = handler(request)
            logger.info(f"Task {request.task_id} completed with status: {response.status.value}")
            return response
        except Exception as e:
            logger.exception(f"Task {request.task_id} failed with error: {e}")
            return TaskResponse.failure(request.task_id, str(e), request.context)
    
    async def run(self) -> NoReturn:
        """Main worker loop."""
        await self.connect()
        self.running = True
        
        logger.info(f"Worker {WORKER_ID} started")
        logger.info(f"Registered handlers: {list_handlers()}")
        
        if self._redis is None:
            logger.info("Running in standalone mode - no task queue processing")
            logger.info("Worker is ready to process tasks via direct API calls")
            while self.running:
                await asyncio.sleep(1)
            return
        
        logger.info(f"Listening on queue: {self.task_queue}")
        
        while self.running:
            try:
                result = await self._redis.blpop(self.task_queue, timeout=1)
                if result is None:
                    continue
                
                _, task_json = result
                request = TaskRequest.from_json(task_json.decode())
                response = await self.process_task(request)
                
                await self._redis.rpush(self.result_queue, response.to_json())
                
            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.exception(f"Error processing task: {e}")
                await asyncio.sleep(1)
        
        await self.disconnect()
    
    def stop(self) -> None:
        """Stop the worker."""
        logger.info("Stopping worker...")
        self.running = False


async def main() -> None:
    """Main entry point."""
    worker = Worker(REDIS_URL, TASK_QUEUE, RESULT_QUEUE)
    
    loop = asyncio.get_event_loop()
    
    def signal_handler() -> None:
        worker.stop()
    
    for sig in (signal.SIGTERM, signal.SIGINT):
        loop.add_signal_handler(sig, signal_handler)
    
    try:
        await worker.run()
    except KeyboardInterrupt:
        worker.stop()


if __name__ == "__main__":
    asyncio.run(main())

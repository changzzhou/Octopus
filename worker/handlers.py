"""
Task Handlers for Octopus Worker

This module contains the task handler implementations. Each handler processes
a specific task_type and returns a TaskResponse.
"""

import logging
from typing import Callable
from protocol import TaskRequest, TaskResponse

logger = logging.getLogger(__name__)

TaskHandler = Callable[[TaskRequest], TaskResponse]

_handlers: dict[str, TaskHandler] = {}


def register_handler(task_type: str) -> Callable[[TaskHandler], TaskHandler]:
    """Decorator to register a task handler for a specific task type."""
    def decorator(handler: TaskHandler) -> TaskHandler:
        _handlers[task_type] = handler
        logger.info(f"Registered handler for task type: {task_type}")
        return handler
    return decorator


def get_handler(task_type: str) -> TaskHandler | None:
    """Get the handler for a specific task type."""
    return _handlers.get(task_type)


def list_handlers() -> list[str]:
    """List all registered task types."""
    return list(_handlers.keys())


@register_handler("noop")
def handle_noop(request: TaskRequest) -> TaskResponse:
    """No-op handler for testing."""
    logger.info(f"Executing noop task: {request.task_id}")
    return TaskResponse.success(request.task_id, {"message": "noop completed"}, request.context)


@register_handler("echo")
def handle_echo(request: TaskRequest) -> TaskResponse:
    """Echo handler - returns the payload as result."""
    logger.info(f"Executing echo task: {request.task_id}")
    return TaskResponse.success(request.task_id, {"echo": request.payload}, request.context)


@register_handler("http")
def handle_http(request: TaskRequest) -> TaskResponse:
    """
    HTTP request handler stub.
    
    Expected payload:
    {
        "url": "https://example.com/api",
        "method": "GET",
        "headers": {},
        "body": null
    }
    
    Full implementation will be added in W-1.
    """
    logger.info(f"Executing http task: {request.task_id}")
    payload = request.payload or {}
    url = payload.get("url", "") if isinstance(payload, dict) else ""
    method = payload.get("method", "GET") if isinstance(payload, dict) else "GET"
    return TaskResponse.success(
        request.task_id,
        {
            "stub": True,
            "message": f"HTTP {method} to {url} - stub implementation",
        },
        request.context,
    )


@register_handler("script")
def handle_script(request: TaskRequest) -> TaskResponse:
    """
    Script execution handler stub.
    
    Expected payload:
    {
        "language": "python",
        "code": "print('hello')",
        "timeout": 30
    }
    
    Full implementation will be added in W-1.
    """
    logger.info(f"Executing script task: {request.task_id}")
    payload = request.payload or {}
    language = payload.get("language", "python") if isinstance(payload, dict) else "python"
    return TaskResponse.success(
        request.task_id,
        {
            "stub": True,
            "message": f"Script ({language}) execution - stub implementation",
        },
        request.context,
    )

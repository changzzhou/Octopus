"""
Unified Task Protocol for Octopus Workers

This module defines the standard task request/response protocol that enables
multi-language worker implementations. All workers (Python, Go, Node.js, etc.)
should follow this protocol for interoperability.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any
import json
import uuid
from datetime import datetime


class TaskStatus(Enum):
    """Task execution status outcomes."""
    PENDING = "pending"
    SUCCESS = "success"
    FAILURE = "failure"


@dataclass
class TaskContext:
    """Execution context for a task within a workflow."""
    workflow_id: str
    execution_id: str
    node_id: str
    attempt: int = 1
    
    def to_dict(self) -> dict[str, Any]:
        return {
            "workflow_id": self.workflow_id,
            "execution_id": self.execution_id,
            "node_id": self.node_id,
            "attempt": self.attempt,
        }
    
    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "TaskContext":
        return cls(
            workflow_id=data["workflow_id"],
            execution_id=data["execution_id"],
            node_id=data["node_id"],
            attempt=data.get("attempt", 1),
        )


@dataclass
class TaskRequest:
    """
    Standard task request format.
    
    This is the payload sent from the workflow engine to workers.
    """
    task_id: str
    task_type: str
    payload: dict[str, Any]
    context: TaskContext
    created_at: str = field(default_factory=lambda: datetime.utcnow().isoformat())
    
    def to_dict(self) -> dict[str, Any]:
        return {
            "task_id": self.task_id,
            "task_type": self.task_type,
            "payload": self.payload,
            "context": self.context.to_dict(),
            "created_at": self.created_at,
        }
    
    def to_json(self) -> str:
        return json.dumps(self.to_dict())
    
    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "TaskRequest":
        return cls(
            task_id=data["task_id"],
            task_type=data["task_type"],
            payload=data["payload"],
            context=TaskContext.from_dict(data["context"]),
            created_at=data.get("created_at", datetime.utcnow().isoformat()),
        )
    
    @classmethod
    def from_json(cls, json_str: str) -> "TaskRequest":
        return cls.from_dict(json.loads(json_str))


@dataclass
class TaskResponse:
    """
    Standard task response format.
    
    This is the payload sent from workers back to the workflow engine.
    """
    task_id: str
    status: TaskStatus
    result: dict[str, Any] | None = None
    error: str | None = None
    completed_at: str = field(default_factory=lambda: datetime.utcnow().isoformat())
    
    def to_dict(self) -> dict[str, Any]:
        return {
            "task_id": self.task_id,
            "status": self.status.value,
            "result": self.result,
            "error": self.error,
            "completed_at": self.completed_at,
        }
    
    def to_json(self) -> str:
        return json.dumps(self.to_dict())
    
    @classmethod
    def success(cls, task_id: str, result: dict[str, Any] | None = None) -> "TaskResponse":
        """Create a successful task response."""
        return cls(task_id=task_id, status=TaskStatus.SUCCESS, result=result)
    
    @classmethod
    def failure(cls, task_id: str, error: str) -> "TaskResponse":
        """Create a failed task response."""
        return cls(task_id=task_id, status=TaskStatus.FAILURE, error=error)
    
    @classmethod
    def pending(cls, task_id: str) -> "TaskResponse":
        """Create a pending task response (for async tasks)."""
        return cls(task_id=task_id, status=TaskStatus.PENDING)


def create_task_request(
    task_type: str,
    payload: dict[str, Any],
    workflow_id: str,
    execution_id: str,
    node_id: str,
) -> TaskRequest:
    """Helper function to create a new task request with auto-generated task_id."""
    return TaskRequest(
        task_id=str(uuid.uuid4()),
        task_type=task_type,
        payload=payload,
        context=TaskContext(
            workflow_id=workflow_id,
            execution_id=execution_id,
            node_id=node_id,
        ),
    )

"""Minimal Python SDK client for the ACF sidecar."""

from __future__ import annotations

from typing import Any, Dict

import requests


class Firewall:
    """Policy Enforcement Point client that forwards input to the sidecar."""

    def __init__(
        self,
        base_url: str = "http://localhost:8080",
        timeout: float = 5.0,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout

    def send(self, payload: str) -> Dict[str, Any]:
        """Send one payload to the sidecar and return the JSON decision."""
        response = requests.post(
            f"{self.base_url}/evaluate",
            json={"input": payload},
            timeout=self.timeout,
        )
        response.raise_for_status()
        return response.json()

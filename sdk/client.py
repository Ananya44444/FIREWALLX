"""Small helper entry point for sending a single payload."""

from __future__ import annotations

import argparse
import json

from firewall import Firewall


def main() -> None:
    parser = argparse.ArgumentParser(description="Send one payload to the ACF sidecar")
    parser.add_argument("payload", help="Input payload to evaluate")
    parser.add_argument("--url", default="http://localhost:8080", help="Sidecar base URL")
    args = parser.parse_args()

    fw = Firewall(base_url=args.url)
    result = fw.send(args.payload)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()

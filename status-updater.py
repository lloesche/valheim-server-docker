#!/usr/bin/env python3
#
#    Copyright 2021 Lukas Lösche
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
import os
import a2s
import json
import socket
import time
import logging
from typing import Dict
from signal import signal, SIGTERM, SIGINT
from argparse import ArgumentParser
from datetime import datetime, timezone
from pprint import pformat

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)
log = logging.getLogger("status-updater")
log.setLevel(logging.INFO)

run = True


def main() -> None:
    parser = get_arg_parser()
    args = parser.parse_args()
    if args.verbose:
        log.setLevel(logging.DEBUG)

    query_host = args.host
    query_port = args.port + 1
    status_file = args.status_file
    signal(SIGINT, handler)
    signal(SIGTERM, handler)
    log.info("Valheim status updater started")
    while run:
        status = get_status(query_host, query_port)
        write_status(status, status_file)
        time.sleep(10)


def write_status(status: str, status_file: str) -> bool:
    tmp_status_file = status_file + ".tmp"
    json_output = json.dumps(status, default=json_default, skipkeys=True)
    with open(tmp_status_file, "w") as json_file:
        json_file.write(json_output)
    os.replace(tmp_status_file, status_file)


def get_status(query_host: str, query_port: int) -> Dict:
    status = {
        "last_status_update": datetime.utcnow().replace(tzinfo=timezone.utc),
        "error": None,
    }
    try:
        info = a2s.info((query_host, query_port))
        players = a2s.players((query_host, query_port))
    except Exception as e:
        status.update({"error": e})
    else:
        status.update(
            {
                "server_name": info.server_name,
                "server_type": info.server_type,
                "platform": info.platform,
                "player_count": len(players),
                "password_protected": info.password_protected,
                "vac_enabled": info.vac_enabled,
                "port": info.port,
                "steam_id": info.steam_id,
                "keywords": info.keywords,
                "game_id": info.game_id,
                "players": [
                    {"name": pl.name, "score": pl.score, "duration": pl.duration}
                    for pl in players
                ],
            }
        )
    return status


def handler(sig, frame) -> None:
    global run
    run = False


def json_default(o):
    if hasattr(o, "to_json"):
        return o.to_json()
    elif isinstance(o, datetime):
        return o.isoformat()
    elif isinstance(o, Exception):
        return pformat(o)
    raise TypeError(f"Object of type {o.__class__.__name__} is not JSON serializable")


def get_arg_parser() -> ArgumentParser:
    parser = ArgumentParser(description="Valheim status updater")
    parser.add_argument(
        "--verbose",
        "-v",
        help="Verbose logging",
        dest="verbose",
        action="store_true",
        default=False,
    )
    parser.add_argument(
        "--host",
        help="Valheim Server Hostname (default: localhost)",
        dest="host",
        type=str,
        default="localhost",
    )
    default_port = int(os.getenv("SERVER_PORT", 2456))
    parser.add_argument(
        "--port",
        help=f"Valheim Server Port (default: {default_port})",
        dest="port",
        type=int,
        default=default_port,
    )
    default_htdocs = os.getenv("STATUS_HTTP_HTDOCS", "/opt/valheim/htdocs")
    default_status_file = f"{default_htdocs}/status.json"
    parser.add_argument(
        "--status-file",
        help=f"Server status file (default: {default_status_file})",
        dest="status_file",
        type=str,
        default=default_status_file,
    )

    return parser


if __name__ == "__main__":
    main()

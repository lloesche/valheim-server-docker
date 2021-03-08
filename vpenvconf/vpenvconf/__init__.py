import os
import logging
from configparser import ConfigParser
from argparse import ArgumentParser
from collections import defaultdict
from typing import Dict

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)
log = logging.getLogger("vpenvconf")
log.setLevel(logging.INFO)


def write_config(config: ConfigParser, config_file: str) -> None:

    config_directory = os.path.dirname(config_file)
    if not os.path.isdir(config_directory):
        log.info(
            f"ValheimPlus config directory {config_directory} does not exist - creating"
        )
        os.makedirs(config_directory)

    new_config_file = config_file + ".tmp"
    log.info(f"Writing ValheimPlus config {new_config_file}")
    with open(new_config_file, "w") as file:
        config.write(file, space_around_delimiters=False)

    if os.path.isfile(config_file):
        old_config_file = config_file + ".old"
        log.info(f"Moving old config {config_file} -> {old_config_file}")
        os.replace(config_file, old_config_file)

    log.info(f"Moving new config {new_config_file} -> {config_file}")
    os.replace(new_config_file, config_file)


def merge_env_with_config(config: ConfigParser, env: Dict) -> None:
    for section, section_content in env.items():
        log.debug(f"Processing section {section}")
        for key, value in section_content.items():
            log.debug(f"Setting {key} = {value} in config section {section}")
            if section not in config:
                config[section] = {}
            config[section][key] = value


def get_env(prefix: str) -> Dict:
    log.debug(f"Reading ValheimPlus config from env variables prefixed with {prefix}")
    env_config = defaultdict(dict)
    env = {k[len(prefix) :]: v for k, v in os.environ.items() if k.startswith(prefix)}
    for k, v in env.items():
        section, key = k.split("_", 1)
        env_config[section][key] = v
        log.debug(f"Found [{section}] {key} = {v} in environment")
    return dict(env_config)


def get_config(config_file: str) -> ConfigParser:
    config = ConfigParser()
    config.optionxform = str
    if os.path.isfile(config_file):
        log.debug(f"Reading existing ValheimPlus config {config_file}")
        config.read(config_file)
    return config


def get_arg_parser() -> ArgumentParser:
    parser = ArgumentParser(description="Generate ValheimPlus config from env")
    parser.add_argument(
        "--verbose",
        "-v",
        help="Verbose logging",
        dest="verbose",
        action="store_true",
        default=False,
    )
    return parser


def add_args(parser: ArgumentParser) -> None:
    parser.add_argument(
        "--config",
        help=(
            "Path to ValheimPlus config file"
            " (default: /config/valheimplus/valheim_plus.cfg)"
        ),
        dest="config_file",
        type=str,
        default="/config/valheimplus/valheim_plus.cfg",
    )
    parser.add_argument(
        "--env-prefix",
        help="Environment prefix (default: VPCFG_)",
        dest="env_prefix",
        type=str,
        default="VPCFG_",
    )

import os
import logging
from configparser import ConfigParser
from argparse import ArgumentParser
from collections import defaultdict
from typing import Dict, Tuple

logging.basicConfig(
    level=logging.INFO,
    format="%(levelname)s - %(message)s",
)
log = logging.getLogger("env2cfg")
log.setLevel(logging.INFO)


_PRE_TRANSLATE = {
    "_DOT_": ".",
    "_HYPHEN_": "-",
    "_PLUS_": "+",
    "_UNDERSCORE_": "@UNDERSCORE@",
}

_POST_TRANSLATE = {"@UNDERSCORE@": "_"}


def write_config(config: ConfigParser, config_file: str) -> None:

    config_directory = os.path.dirname(config_file)
    if not os.path.isdir(config_directory):
        log.info(f"Mod config directory {config_directory} does not exist - creating")
        os.makedirs(config_directory)

    new_config_file = config_file + ".tmp"
    log.info(f"Writing mod config {new_config_file}")
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
    log.debug(f"Reading mod config from env variables prefixed with {prefix}")
    env_config = defaultdict(dict)
    env = {k[len(prefix) :]: v for k, v in os.environ.items() if k.startswith(prefix)}
    for k, v in env.items():
        section, key = var_process(k)
        env_config[section][key] = v
        log.debug(f"Found [{section}] {key} = {v} in environment")
    return dict(env_config)


def var_process(var: str) -> Tuple:
    for in_char, out_char in _PRE_TRANSLATE.items():
        var = var.replace(in_char, out_char)
    section, key = var.split("_", 1)
    for in_char, out_char in _POST_TRANSLATE.items():
        section = section.replace(in_char, out_char)
        key = key.replace(in_char, out_char)
    return section, key


def get_config(config_file: str) -> ConfigParser:
    config = ConfigParser()
    config.optionxform = str
    if os.path.isfile(config_file):
        log.debug(f"Reading existing mod config {config_file}")
        with open(config_file, mode="rb") as f:
            content = f.read()
            if content.startswith(b"\xef\xbb\xbf"):
                config.read_string(content.decode("utf-8-sig"))
            else:
                config.read_string(content.decode("utf-8"))
    return config


def get_arg_parser() -> ArgumentParser:
    parser = ArgumentParser(description="Generate mod config from env")
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
        help=("Path to mod config file" " (default: /config/bepinex/BepInEx.cfg)"),
        dest="config_file",
        type=str,
        default="/config/bepinex/BepInEx.cfg",
    )
    parser.add_argument(
        "--env-prefix",
        help="Environment prefix (default: MODCFG_)",
        dest="env_prefix",
        type=str,
        default="MODCFG_",
    )

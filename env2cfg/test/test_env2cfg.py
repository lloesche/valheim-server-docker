import os
import tempfile
import logging
from configparser import ConfigParser
from env2cfg import (
    log,
    get_arg_parser,
    add_args,
    get_env,
    get_config,
    merge_env_with_config,
    write_config,
)

log.setLevel(logging.DEBUG)
parser = get_arg_parser()
add_args(parser)
args = parser.parse_args()

test_var1 = args.env_prefix + "Server_test"
test_var2 = args.env_prefix + "Logging_DOT_Console_Enabled"
test_var3 = args.env_prefix + "Logging_DOT_Console_Some_UNDERSCORE_Var"
os.environ[test_var1] = "true"
os.environ[test_var2] = "true"
os.environ[test_var3] = "false"


def test_args():
    assert args.verbose is False
    assert args.config_file == "/config/bepinex/BepInEx.cfg"
    assert args.env_prefix == "MODCFG_"


def test_get_env():
    env = get_env(args.env_prefix)
    assert "Server" in env
    assert "test" in env["Server"]
    assert env["Server"]["test"] == "true"
    assert env["Logging.Console"]["Enabled"] == "true"
    assert env["Logging.Console"]["Some_Var"] == "false"


def test_merge_env_with_config():
    env = get_env(args.env_prefix)
    config = ConfigParser()
    config.optionxform = str
    merge_env_with_config(config, env)
    assert "Server" in config
    assert "test" in config["Server"]
    assert config["Server"]["test"] == "true"


def test_read_write_config():
    temp_configfile = tempfile.NamedTemporaryFile(delete=False)
    config_file = temp_configfile.name
    temp_configfile.close()

    config = ConfigParser()
    config.optionxform = str
    config["Server"] = {}
    config["Server"]["test"] = "true"
    write_config(config, config_file)

    new_config = get_config(config_file)
    os.unlink(config_file)

    assert "Server" in new_config
    assert "test" in new_config["Server"]
    assert new_config["Server"]["test"] == "true"

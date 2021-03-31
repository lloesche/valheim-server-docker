import logging
from . import (
    log,
    get_arg_parser,
    add_args,
    get_env,
    get_config,
    merge_env_with_config,
    write_config,
)


def main() -> None:
    parser = get_arg_parser()
    add_args(parser)
    args = parser.parse_args()
    if args.verbose:
        log.setLevel(logging.DEBUG)

    env = get_env(args.env_prefix)
    if len(env) > 0:
        config = get_config(args.config_file)
        merge_env_with_config(config, env)
        write_config(config, args.config_file)
    else:
        log.info("No BepInEx config found in env")


if __name__ == "__main__":
    main()

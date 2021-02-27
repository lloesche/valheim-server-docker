from vpenvconf import get_arg_parser, add_args


def test_args():
    parser = get_arg_parser()
    add_args(parser)
    args = parser.parse_args()
    assert args.verbose is False
    assert args.config_file == "/config/valheimplus/valheim_plus.cfg"
    assert args.env_prefix == "VPCFG_"

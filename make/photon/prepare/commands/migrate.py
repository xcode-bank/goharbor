import os, sys, importlib, shutil, glob

import click

from utils.misc import get_realpath
from migration.utils import read_conf, search


@click.command()
@click.option('-i', '--input', 'input_', default='', help="The path of original config file")
@click.option('-o', '--output', default='', help="the path of output config file")
@click.option('-t', '--target', default='2.0.0', help="target version of input path")
def migrate(input_, output, target):
    if not output:
        output = input_
    input_path = get_realpath(input_)
    output_path = get_realpath(output)

    configs = read_conf(input_path)

    input_version = configs.get('_version')
    
    if input_version == target:
        click.echo("Version of input harbor.yml is identical to target {}, no need to upgrade".format(input_version))
        sys.exit(0)

    current_input_path = input_path
    for m in search(input_version, target):
        current_output_path = "harbor.yml.{}.tmp".format(m.revision)
        click.echo("migrating to version {}".format(m.revision))
        m.migrate(current_input_path, current_output_path)
        current_input_path = current_output_path
    shutil.copy(current_input_path, output_path)
    click.echo("Written new values to {}".format(output))
    for tmp_f in glob.glob("harbor.yml.*.tmp"):
        os.remove(tmp_f)


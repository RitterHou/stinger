# -*- coding: utf-8 -*-
import os

filename = 'version.txt'


def run():
    with open(filename) as f:
        old_version = f.readline()

    v = old_version.split('.')
    v[2] = str(int(v[2]) + 1)
    new_version = '.'.join(v)

    with open(filename, 'w') as f:
        f.write(new_version)

    os.system('git add version.txt')
    os.system('git commit -m "Version {}"'.format(new_version))
    os.system('git tag {}'.format(new_version))
    os.system('git push origin {}'.format(new_version))

    print("Changed version from {} to {}".format(old_version, new_version))


if __name__ == '__main__':
    run()

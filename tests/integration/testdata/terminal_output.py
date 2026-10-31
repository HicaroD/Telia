"""Check native output against protocol bytes, not against its own encoder."""
import os
import pty
import select
import subprocess
import sys

master, slave = pty.openpty()
try:
    env = dict(os.environ, TERM="xterm-256color", COLORTERM="truecolor")
    child = subprocess.Popen([sys.argv[1]], stdin=slave, stdout=slave, stderr=subprocess.PIPE, env=env)
    try:
        expected = b"\x1b[38;2;12;34;56m\x1b[0m"
        actual = b""
        while len(actual) < len(expected):
            assert select.select([master], [], [], 3)[0], ("missing output", actual)
            actual += os.read(master, 1024)
        _, err = child.communicate(timeout=3)
        assert child.returncode == 0, (child.returncode, err)
        assert actual == expected, (actual, expected)
        assert not select.select([master], [], [], 0)[0], "unexpected extra terminal output"
    finally:
        if child.poll() is None:
            child.kill()
            child.wait()
finally:
    os.close(master)
    os.close(slave)

"""Exercise compiled Telia through a real PTY, without changing the host terminal."""
import fcntl
import os
import pty
import select
import struct
import subprocess
import sys
import termios

master, slave = pty.openpty()
mode = sys.argv[2]

def controlling_terminal():
    os.setsid()
    fcntl.ioctl(slave, termios.TIOCSCTTY, 0)

try:
    fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", 24, 80, 0, 0))
    before = termios.tcgetattr(slave)
    child = subprocess.Popen([sys.argv[1]], stdin=slave, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                             preexec_fn=controlling_terminal if mode == "controlling" else None)
    try:
        assert select.select([child.stdout], [], [], 3)[0], "child did not become ready"
        assert child.stdout.read(1) == b"R", "child exited before entering raw mode"
        if mode not in ("controlling", "cleanup"):
            during = termios.tcgetattr(slave)
            assert not during[3] & (termios.ECHO | termios.ICANON | termios.ISIG)
            if mode == "resize":
                fcntl.ioctl(slave, termios.TIOCSWINSZ, struct.pack("HHHH", 40, 100, 0, 0))
                assert select.select([child.stdout], [], [], 3)[0], "child missed resize"
                assert child.stdout.read(1) == b"R", "child did not report resize"
            os.write(master, b"\x03")
        out, err = child.communicate(timeout=3)
        assert child.returncode == 0, (child.returncode, out, err)
        if mode != "controlling":
            after = termios.tcgetattr(slave)
            # macOS sets kernel-maintained PENDIN when restoring canonical input.
            pending = getattr(termios, "PENDIN", 0)
            before[3] &= ~pending
            after[3] &= ~pending
            assert after == before, ("native terminal state was not restored exactly", before, after)
    finally:
        if child.poll() is None:
            child.kill()
            child.wait()
finally:
    os.close(master)
    os.close(slave)

"""Run one synthetic Cerebras completion through the local A4 stack."""

import argparse
import json
import os
from pathlib import Path
from queue import Empty, Queue
from shlex import quote
import signal
import subprocess
from tempfile import TemporaryDirectory
from threading import Thread


def read_lab(output, result):
    fixture = {}
    for line in output:
        if "Cerebras smoke " in line:
            print(line.rstrip(), flush=True)
        for prefix, name in (
            ("Open ", "origin"),
            ("Synthetic president student ID: ", "login"),
            ("Synthetic password for both users: ", "password"),
        ):
            if line.startswith(prefix):
                value = line.removeprefix(prefix).strip()
                fixture[name] = value.split("/#", 1)[0] if name == "origin" else value
        if line.startswith("Stop with Ctrl+C"):
            result.put(fixture)
    result.put(None)


def run(args):
    smoke = Path(__file__).with_name("cerebras-smoke-h5.mjs").resolve()
    env = dict(os.environ, CAMS_ENV_FILE=os.devnull)
    env.pop("CEREBRAS_API_KEY", None)

    with TemporaryDirectory(prefix="campusly-cerebras-lab-") as directory:
        wrapper = Path(directory) / "assistant"
        wrapper.write_text(
            "#!/bin/sh\nexec " + quote(str(args.binaries / "assistant"))
            + " -cerebras-smoke-key-file " + quote(str(args.key_file)) + "\n"
        )
        wrapper.chmod(0o700)
        (Path(directory) / "fixture").symlink_to(args.binaries / "fixture")
        command = [
            args.core_python, "-m", "scripts.m2.assistant_lab", directory,
            str(args.h5), "--review",
        ]
        with subprocess.Popen(
            command, cwd=args.core, env=env, stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT, text=True, start_new_session=True,
        ) as lab:
            ready = Queue()
            reader = Thread(target=read_lab, args=(lab.stdout, ready), daemon=True)
            reader.start()
            try:
                try:
                    fixture = ready.get(timeout=180)
                except Empty:
                    raise RuntimeError("The disposable A4 lab did not become ready") from None
                if not fixture or set(fixture) != {"origin", "login", "password"}:
                    raise RuntimeError("The disposable A4 lab exited before providing synthetic credentials")
                fixture["screenshot"] = str(args.screenshot)
                print("Disposable A4 lab ready. Starting one browser scenario.", flush=True)
                subprocess.run(
                    ["node", smoke], check=True, input=json.dumps(fixture), text=True, env=env,
                    timeout=240,
                )
            finally:
                if lab.poll() is None:
                    lab.send_signal(signal.SIGINT)
                lab.wait(timeout=45)
                reader.join(timeout=2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--core", type=Path, required=True)
    parser.add_argument("--h5", type=Path, required=True)
    parser.add_argument("--binaries", type=Path, required=True)
    parser.add_argument("--core-python", required=True)
    parser.add_argument("--key-file", type=Path, required=True)
    parser.add_argument("--screenshot", type=Path, required=True)
    args = parser.parse_args()
    for field in ("core", "h5", "binaries", "key_file", "screenshot"):
        setattr(args, field, getattr(args, field).absolute())
    run(args)

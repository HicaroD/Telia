ALGORITHMS = [
    "fib"
]


def compile(path: str):
    # TODO:
    pass


def bench(path: str):
    pass


def main():
    for algorithm in ALGORITHMS:
        path = f"code/{algorithm}"
        compile(path)
        bench(path)


if __name__ == "__main__":
    main()

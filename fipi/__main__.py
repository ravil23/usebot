import argparse

from fipi.crawler import FIPICrawler


def main(cache_dir: str, output_dir: str, session_id: str, force: bool) -> None:
    crawler = FIPICrawler(cache_dir, output_dir, session_id, force)
    crawler.load_dictionaries()
    tasks_subject_russian = crawler.load_subject_russian()
    # TODO: remove tasks_subject_russian filters
    crawler.save_subject_russian([task for task in tasks_subject_russian if task.type_id == 2 and task.doc is None])


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Load FIPI tasks.')
    parser.add_argument('--cache', type=str, required=True, help='cache directory')
    parser.add_argument('--output', type=str, required=True, help='output directory')
    parser.add_argument('--session', type=str, required=False, help='session id')
    parser.add_argument('--force', action='store_true', help='overwrite existed data')

    args = parser.parse_args()

    main(args.cache, args.output, args.session, args.force)

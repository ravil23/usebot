import argparse

from crawler.fipi import FIPICrawler

SITE_FIPI = "fipi"


def main(cache_dir: str, output_dir: str, site: str, fipi_session_id: str, force: bool) -> None:
    crawler = None
    if site == SITE_FIPI:
        crawler = FIPICrawler(cache_dir, output_dir, fipi_session_id, force)
        crawler.load_dictionaries()
        tasks_subject_russian = crawler.load_subject_russian()
        # TODO: remove tasks_subject_russian filters
        crawler.save_subject_russian([task for task in tasks_subject_russian if task.type_id == 2 and task.doc is None])
    else:
        raise IllegalArgumentError()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Load tasks.')
    parser.add_argument('--cache', type=str, required=True, help='cache directory')
    parser.add_argument('--output', type=str, required=True, help='output directory')
    parser.add_argument('--site', type=str, required=True, choices=[SITE_FIPI], help='site for crawling')
    parser.add_argument('--fipi-session', type=str, required=False, help='FIPI session id')
    parser.add_argument('--force', action='store_true', help='overwrite existed data')

    args = parser.parse_args()

    main(args.cache, args.output, args.site, args.fipi_session, args.force)

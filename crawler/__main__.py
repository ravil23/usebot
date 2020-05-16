import argparse

from crawler.fipi import FIPICrawler

SITE_FIPI = "fipi"


def main(cache_dir: str, output_dir: str, site: str, fipi_session_id: str, force: bool) -> None:
    if site == SITE_FIPI:
        crawler = FIPICrawler(cache_dir, output_dir, fipi_session_id, force)
        crawler.load_dictionaries()

        subjects = crawler.load_subjects()
        for subject_id, tasks in subjects.items():
            crawler.save_subject(
                [
                    task
                    for task in tasks
                    if task.type_id == 2
                    and len(task.text) <= 500
                    and '=' not in task.text
                    and len(task.options) == 4
                    and not (task.subject_id == 3 and any(requirement.startswith("2.4 ") for requirement in task.requirements))
                ],
                subject_id,
                crawler.SUBJECT_FILENAMES[subject_id]
            )
    else:
        raise ValueError(f'invalid site: {site}')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Load tasks.')
    parser.add_argument('--cache', type=str, required=True, help='cache directory')
    parser.add_argument('--output', type=str, required=True, help='output directory')
    parser.add_argument('--site', type=str, required=True, choices=[SITE_FIPI], help='site for crawling')
    parser.add_argument('--fipi-session', type=str, required=False, help='FIPI session id')
    parser.add_argument('--force', action='store_true', help='overwrite existed data')

    args = parser.parse_args()

    main(args.cache, args.output, args.site, args.fipi_session, args.force)

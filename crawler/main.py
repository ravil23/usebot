import argparse
import json
import logging
import os

import requests


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class FIPICrawler:
    API_DICTIONARIES = 'http://os.fipi.ru/api/dictionaries'
    API_TASKS = 'http://os.fipi.ru/api/tasks'

    DICTIONARIES_FILENAME = 'dictionaries.json'
    TASKS_SUBJECT_RUSSIAN_FILENAME = 'tasks_subject_russian.json'

    def __init__(self, output_dir: str, session_id: str) -> None:
        self.output_dir = output_dir
        self.session_id = session_id

    def load_dictionaries(self) -> None:
        response = requests.get(self.API_DICTIONARIES)
        response.raise_for_status()
        output_path = os.path.join(self.output_dir, self.DICTIONARIES_FILENAME)
        self._dump(response.json(), output_path)
        logging.info('Dictionaries loaded')

    def load_subject_russian(self) -> None:
        headers = {'sessionId': self.session_id}
        request_data = {
            'subjectId': '1',
            'levelIds': [],
            'themeIds': [],
            'typeIds': [],
            'id': '',
            'favorites': 0,
            'answerStatus': 0,
            'themeSectionIds': [],
            'published': 0,
            'extId': '',
            'fipiCode': '',
            'docId': '',
            'isAdmin': False,
            'loadDates': [],
            'isPublished': False,
            'pageSize': 100,
            'pageNumber': 1,
        }
        tasks = []
        while True:
            response = requests.post(self.API_TASKS, headers=headers, json=request_data)
            response.raise_for_status()
            page_tasks = response.json()['tasks']
            if len(page_tasks) == 0:
                break
            tasks.extend(page_tasks)
            request_data['pageNumber'] += 1
            logging.info(f'{len(tasks)} tasks for russian subject loaded')
        output_path = os.path.join(self.output_dir, self.TASKS_SUBJECT_RUSSIAN_FILENAME)
        self._dump({'tasks': tasks}, output_path)

    def _dump(self, data: dict, output_path: str) -> None:
        with open(output_path, 'w', encoding='utf8') as f:
            json.dump(data, f, indent='  ', ensure_ascii=False)


def main(output_dir: str, session_id: str) -> None:
    crawler = FIPICrawler(output_dir, session_id)
    crawler.load_dictionaries()
    crawler.load_subject_russian()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process some integers.')
    parser.add_argument('--out', type=str, required=True, help='output directory')
    parser.add_argument('--session', type=str, required=True, help='session id')

    args = parser.parse_args()

    main(args.out, args.session)

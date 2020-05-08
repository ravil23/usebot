import json
import logging
import os
from typing import Dict, List

import requests

from fipi.task import FIPITask

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


class FIPICrawler:
    API_DICTIONARIES = 'http://os.fipi.ru/api/dictionaries'
    API_TASKS = 'http://os.fipi.ru/api/tasks'

    DICTIONARIES_FILENAME = 'dictionaries.json'
    TASKS_SUBJECT_RUSSIAN_FILENAME = 'tasks_subject_russian.json'

    def __init__(self, cache_dir: str, output_dir: str, session_id: str, force: bool) -> None:
        self.cache_dir = cache_dir
        self.output_dir = output_dir
        self.session_id = session_id
        self.force = force

    def load_dictionaries(self) -> Dict:
        output_path = os.path.join(self.cache_dir, self.DICTIONARIES_FILENAME)
        if not self.force and os.path.exists(output_path):
            with open(output_path, 'r') as f:
                dictionaries = json.load(f)
            logging.info('Dictionaries loaded from cache')
            return dictionaries

        response = requests.get(self.API_DICTIONARIES)
        response.raise_for_status()
        dictionaries = response.json()

        self._dump(dictionaries, output_path)
        logging.info('Dictionaries loaded from site')
        return dictionaries

    def load_subject_russian(self) -> List[FIPITask]:
        output_path = os.path.join(self.cache_dir, self.TASKS_SUBJECT_RUSSIAN_FILENAME)
        if not self.force and os.path.exists(output_path):
            with open(output_path, 'r') as f:
                tasks = [FIPITask.from_response(raw_task) for raw_task in json.load(f)['tasks']]
            logging.info(f'{len(tasks)} tasks for russian subject loaded from cache')
            return tasks

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
            logging.info(f'{len(tasks)} tasks for russian subject loaded from site')

        self._dump({'tasks': tasks}, output_path)
        return [FIPITask.from_response(raw_task) for raw_task in tasks ]

    def save_subject_russian(self, tasks: List[FIPITask]) -> None:
        data = {'tasks': [task.to_dict() for task in tasks]}
        output_path = os.path.join(self.output_dir, self.TASKS_SUBJECT_RUSSIAN_FILENAME)
        self._dump(data, output_path)

    @staticmethod
    def _dump(data: dict, output_path: str) -> None:
        with open(output_path, 'w', encoding='utf8') as f:
            json.dump(data, f, indent='  ', ensure_ascii=False)

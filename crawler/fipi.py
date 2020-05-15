import json
import logging
import os
from typing import Dict, List

import requests

from crawler.task import Task

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')


class FIPICrawler:
    API_DICTIONARIES = 'http://os.fipi.ru/api/dictionaries'
    API_TASKS = 'http://os.fipi.ru/api/tasks'

    DICTIONARIES_FILENAME = 'dictionaries.json'
    TASKS_SUBJECT_RUSSIAN_FILENAME = 'tasks_subject_russian.json'
    TASKS_SUBJECT_MATCH_ADVANCED_FILENAME = 'tasks_subject_math_advanced.json'
    TASKS_SUBJECT_PHYSICS_FILENAME = 'tasks_subject_physics.json'
    TASKS_SUBJECT_CHEMISTRY_FILENAME = 'tasks_subject_chemistry.json'
    TASKS_SUBJECT_IT_FILENAME = 'tasks_subject_it.json'
    TASKS_SUBJECT_BIOLOGY_FILENAME = 'tasks_subject_biology.json'
    TASKS_SUBJECT_HISTORY_FILENAME = 'tasks_subject_history.json'
    TASKS_SUBJECT_GEOGRAPHY_FILENAME = 'tasks_subject_geography.json'
    TASKS_SUBJECT_ENGLISH_FILENAME = 'tasks_subject_english.json'
    TASKS_SUBJECT_GERMAN_FILENAME = 'tasks_subject_german.json'
    TASKS_SUBJECT_FRENCH_FILENAME = 'tasks_subject_french.json'
    TASKS_SUBJECT_SOCIAL_FILENAME = 'tasks_subject_social.json'
    TASKS_SUBJECT_SPANISH_FILENAME = 'tasks_subject_spanish.json'
    TASKS_SUBJECT_LITERATURE_FILENAME = 'tasks_subject_literature.json'
    TASKS_SUBJECT_MATH_BASIC_FILENAME = 'tasks_subject_math_basic.json'

    SUBJECT_ID_RUSSIAN = '1'
    SUBJECT_ID_MATH_ADVANCED = '2'
    SUBJECT_ID_PHYSICS = '3'
    SUBJECT_ID_CHEMISTRY = '4'
    SUBJECT_ID_IT = '5'
    SUBJECT_ID_BIOLOGY = '6'
    SUBJECT_ID_HISTORY = '7'
    SUBJECT_ID_GEOGRAPHY = '8'
    SUBJECT_ID_ENGLISH = '9'
    SUBJECT_ID_GERMAN = '10'
    SUBJECT_ID_FRENCH = '11'
    SUBJECT_ID_SOCIAL = '12'
    SUBJECT_ID_SPANISH = '13'
    SUBJECT_ID_LITERATURE = '18'
    SUBJECT_ID_MATH_BASIC = '22'

    SUBJECT_FILENAMES = {
        SUBJECT_ID_RUSSIAN: TASKS_SUBJECT_RUSSIAN_FILENAME,
        SUBJECT_ID_MATH_ADVANCED: TASKS_SUBJECT_MATCH_ADVANCED_FILENAME,
        SUBJECT_ID_PHYSICS: TASKS_SUBJECT_PHYSICS_FILENAME,
        SUBJECT_ID_CHEMISTRY: TASKS_SUBJECT_CHEMISTRY_FILENAME,
        SUBJECT_ID_IT: TASKS_SUBJECT_IT_FILENAME,
        SUBJECT_ID_BIOLOGY: TASKS_SUBJECT_BIOLOGY_FILENAME,
        SUBJECT_ID_HISTORY: TASKS_SUBJECT_HISTORY_FILENAME,
        SUBJECT_ID_GEOGRAPHY: TASKS_SUBJECT_GEOGRAPHY_FILENAME,
        SUBJECT_ID_ENGLISH: TASKS_SUBJECT_ENGLISH_FILENAME,
        SUBJECT_ID_GERMAN: TASKS_SUBJECT_GERMAN_FILENAME,
        SUBJECT_ID_FRENCH: TASKS_SUBJECT_FRENCH_FILENAME,
        SUBJECT_ID_SOCIAL: TASKS_SUBJECT_SOCIAL_FILENAME,
        SUBJECT_ID_SPANISH: TASKS_SUBJECT_SPANISH_FILENAME,
        SUBJECT_ID_LITERATURE: TASKS_SUBJECT_LITERATURE_FILENAME,
        SUBJECT_ID_MATH_BASIC: TASKS_SUBJECT_MATH_BASIC_FILENAME,
    }

    def __init__(self, cache_dir: str, output_dir: str, session_id: str, force: bool) -> None:
        self.cache_dir = cache_dir
        self.output_dir = output_dir
        self.session_id = session_id
        self.force = force

    def load_dictionaries(self) -> Dict:
        cache_path = os.path.join(self.cache_dir, self.DICTIONARIES_FILENAME)
        if not self.force and os.path.exists(cache_path):
            with open(cache_path, 'r') as f:
                dictionaries = json.load(f)
            logging.info(f'Dictionaries loaded from cache: {cache_path}')
            return dictionaries

        response = requests.get(self.API_DICTIONARIES)
        response.raise_for_status()
        dictionaries = response.json()

        self._dump(dictionaries, cache_path)
        logging.info(f'Dictionaries loaded from site and cached: {cache_path}')
        return dictionaries

    def load_subjects(self) -> Dict[str, List[Task]]:
        return {
            subject_id: self._load_subject(subject_id, filename)
            for subject_id, filename in self.SUBJECT_FILENAMES.items()
        }

    def _load_subject(self, subject_id: str, filename: str) -> List[Task]:
        cache_path = os.path.join(self.cache_dir, filename)
        if not self.force and os.path.exists(cache_path):
            with open(cache_path, 'r') as f:
                tasks = [Task.from_response(raw_task) for raw_task in json.load(f)['tasks']]
            logging.info(f'{len(tasks)} tasks for subject {subject_id} loaded from cache: {cache_path}')
            return tasks

        headers = {'sessionId': self.session_id}
        request_data = {
            'subjectId': subject_id,
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

            try:
                page_tasks = response.json()['tasks']

                if len(page_tasks) == 0:
                    break

                tasks.extend(page_tasks)
                logging.info(f'{len(tasks)} tasks for subject {subject_id} loaded from site')
            except Exception as e:
                logging.error(f'Error on loading subject {subject_id}: {e}')
            finally:
                request_data['pageNumber'] += 1

        self._dump({'tasks': tasks}, cache_path)
        logging.info(f'{len(tasks)} tasks for subject {subject_id} cached: {cache_path}')
        return [Task.from_response(raw_task) for raw_task in tasks ]

    def save_subject(self, tasks: List[Task], subject_id: str, filename: str) -> None:
        data = {'tasks': [task.to_dict() for task in tasks]}
        output_path = os.path.join(self.output_dir, filename)
        self._dump(data, output_path)
        logging.info(f'{len(tasks)} tasks for subject {subject_id} saved: {output_path}')

    @staticmethod
    def _dump(data: dict, output_path: str) -> None:
        with open(output_path, 'w', encoding='utf8') as f:
            json.dump(data, f, indent='  ', ensure_ascii=False)

import logging
import shutil
import tarfile


class Utils:
    def __init__(self) -> None:
        self._logger = logging.getLogger(__name__)

    def to_list(self, value):
        return value if isinstance(value, list) else [value]
    
    def gzip_folder(self, source_folder, output_filename):
        shutil.make_archive(output_filename, 'gztar', source_folder)
    
    def extract_gzip(self, gzip_path, extract_path):
        with tarfile.open(gzip_path, 'r:gz') as tar:
            tar.extractall(path=extract_path)

import os
import re
from typing import Dict, Optional

class TokenParser:
    def __init__(self, config_file: Optional[str] = None):
        self.tokens = {}
        self.config_file = config_file

    # def parse_from_dotenv(self):
    #     with open(self.config_file, 'r') as f:
    #         for line in f.readlines():
    #             if line.startswith('#'):
    #                 continue
    #             elif match := re.search(r'^MULTI_TOKEN(\d+)=(\S+)$', line):
    #                 count = int(match.group(1))
    #                 token = match.group(2)
    #                 if not count and not token:
    #                     continue
    #                 self.tokens[count] = token
    #     return self.config

    def parse_from_env(self) -> Dict[int, str]:
        self.tokens = dict(
            (c+1, t) for c, (_, t) in enumerate(filter(
                lambda n: n[0].startswith('MULTI_TOKEN'),
                    sorted(os.environ.items())
            ))
        )
        return self.tokens
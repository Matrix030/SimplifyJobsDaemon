import argparse
import json
import subprocess
import sys
from pathlib import Path
from xml.sax import parse
from odf.opendocument import load
from odf import text as odf_text

def load_known_projects(json_path: str) -> list[dict]:
    with open(json_path, 'r') as f:
        return json.load(f)

def get_text_content(element) -> str:
    pass

def get_element_tag(element) -> str:
    pass

def is_section_header(element, text: str) -> bool:
    pass

def is_project_title(text: str, known_projects: list[dict]) -> tuple[bool, str | None]:
    pass

def find_projects_section(doc) -> tuple[int | None, int | None]:
    pass

def find_project_blocks(doc, start_idx: int, end_idx: int, known_projects: list[dict]) -> list[dict]:
    pass

def keep_only_projects(doc, projects_to_keep: list[str], known_projects: list[dict]) -> bool:
    pass

def export_to_pdf(odt_path: Path, pdf_path: Path) -> bool:
    pass

def main():
    parser = argparse.ArgumentParser(
        description='Tailor resume by keeping only specified projects'
    )

    parser.add_argument(
        '--template',
        required=True,
        help='Input template ODT file'
    )

    parser.add_argument(
        '--projects-json',
        required=True,
        help='Path to projects.json with known projects'
    )

    parser.add_argument(
        '--keep',
        required=True,
        help='Comma-separated proejct names to keep'
    )

    parser.add_argument(
        '--output',
        required=True,
        help='Output filename (.odt  or .pdf)'
    )

    parser.add_argument(
        '--pdf',
        action='store_true',
        help='Also export to PDF (auto-enabled if output ends with .pdf)'
    )

    args = parser.parse_args()

    

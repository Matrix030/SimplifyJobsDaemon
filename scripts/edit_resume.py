import argparse
import json
import subprocess
import sys
from pathlib import Path
from xml.sax import parse
from odf.opendocument import load
from odf import text as odf_text

def load_known_projects(json_path: str) -> list[dict]:
    #Load project definitions from projects.json
    with open(json_path, 'r') as f:
        return json.load(f)

def get_text_content(element) -> str:
    #Recursively extract text from an ODF element
    result = []
    if hasattr(element, 'childNodes'):
        for node in element.childNodes:
            if node.nodeType == 3:
                result.append(str(node))
            elif hasattr(node, 'childNodes'):
                result.append(get_text_content(node))
    return ''.join(result)

def get_element_tag(element) -> str:
    #Get the tag name of an element
    if hasattr(element, 'qname') and element.qname:
        return element.qname[1] # Returns 'p', 'list', 'h', etc.
    return ''

def is_section_header(element, text: str) -> bool:
    #check if element is a section header
    section_keywords = ['ACHIEVEMENTS', 'CERTIFICATIONS', 'EDUCATION', 'EXPERIENCE', 'SKILLS']
    text_upper = text.upper().strip()

    #check if it's a header-styple element with section keyword
    for keyword in section_keywords:
        if keyword in text_upper and len(text_upper) < 50:
            return True
    return False


def is_project_title(text: str, known_projects: list[dict]) -> tuple[bool, str | None]:
    #Check if text is a project title line
    if "|" not in text:
        return False, None

    #Extract name before "|"
    name_part = text.split("|")[0].strip()

    #Match against known projects
    for proj in known_projects:
        proj_name = proj["name"].lower()
        name_part_lower = name_part.lower()

        #fuzzy match: either contains the other
        if proj_name in name_part_lower or name_part_lower in proj_name:
            return True, proj["name"]

    return False, None


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

    #Validate inputs
    template_path = Path(args.template)
    if not template_path.exists():
        print(f"Error: Template not found: {template_path}")
        sys.exit(1)

    projects_json_path = Path(args.projects_json)
    if not projects_json_path.exists():
        print(f"Error: projects.json not found: {projects_json_path}")
        sys.exit(1)

    #Load known projects
    known_projects = load_known_projects(str(projects_json_path))
    print(f"Loaded {len(known_projects)} known projects from {projects_json_path}")

    #Parse projects to keep
    projects_to_keep = [p.strip() for p in args.keep.split(',')] 
    print(f"Projects to keep: {projects_to_keep}")

    #Load document
    print(f"\nLoading template: {template_path}")
    doc = load(str(template_path))

    #Modify document
    if not keep_only_projects(doc, projects_to_keep, known_projects):
        sys.exit(1)

    #Determine output path
    output_path = Path(args.output)

    if output_path.suffix == ".pdf":
        odt_output = output_path.with_suffix('.odt')
        pdf_output = output_path
        export_pdf = True
    else:
        odt_output = output_path
        pdf_output = output_path.with_suffix('.pdf')
        export_pdf = args.pdf

    #Ensure output directory exists
    odt_output.parent.mkdir(parents=True, exist_ok=True)

    #save modified ODT
    doc.save(str(odt_output))
    print(f"\nSaved ODT: {odt_output}")

    #Export to PDF if requested:
    if export_pdf:
        print(f"Exporting to PDF...")
        if export_to_pdf(odt_output, pdf_output):
            print(f"Saved PDF: {pdf_output}")
        else:
            print("PDF export failed")
            sys.exit(1)

    print("\nDone!")


if __name__ == '__main__':
    main()

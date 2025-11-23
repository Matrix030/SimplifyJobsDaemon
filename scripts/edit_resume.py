import argparse
import json
import subprocess
import sys
from pathlib import Path
from odf.opendocument import load

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
    #find the start and end indices of the PROJECTS section
    body = doc.text
    nodes = list(body.childNodes)

    start_idx = None
    end_idx = None

    for i, elem in enumerate(nodes):
        text = get_text_content(elem).strip()

        #Find PROJECTS header
        if start_idx is None:
            if 'PROJECTS' in text.upper() and len(text) < 30:
                start_idx = i
                continue
        
        #Find end of PROJECTS section
        if start_idx is not None and i > start_idx:
            if is_section_header(elem, text):
                end_idx = i
                break

    #if no end found, assume it goes to the end
    if start_idx is not None and end_idx is None:
        end_idx = len(nodes)

    return start_idx, end_idx


def find_project_blocks(doc, start_idx: int, end_idx: int, known_projects: list[dict]) -> list[dict]:
    #Find all project blocks within the PROJECTS section
    body = doc.text
    nodes = list(body.childNodes)

    projects = []
    current_project = None
    current_start = None

    for i in range(start_idx + 1, end_idx):
        elem = nodes[i]
        tag =get_element_tag(elem)
        text = get_text_content(elem).strip()

        #Check if this is a project title paragraph
        if tag == 'p':
            is_title, proj_name = is_project_title(text, known_projects)

            if is_title:
                #Save previous project block
                if current_project is not None:
                    projects.append({
                                    'name': current_project,
                                    'start_idx': current_start,
                                    'end_idx': i
                                    })

                # Start new project block
                current_project = proj_name
                current_start = i

    #Don't forget the last project
    if current_project is not None:
        projects.append({
                        'name': current_project,
                        'start_idx': current_start,
                        'end_idx': end_idx
        })
    
    return projects


def keep_only_projects(doc, projects_to_keep: list[str], known_projects: list[dict]) -> bool:
    #Remove all projects except the specified ones
    # Find PROJECTS section boundaries
    section_start, section_end = find_projects_section(doc)

    if section_start is None or section_end is None:
        print("Error: Could not find PROJECTS section")
        return False
    
    print(f"Found PROJECTS section: elements {section_start} to {section_end}")

    #Find all project blocks
    project_blocks = find_project_blocks(doc, section_start, section_end, known_projects)
    print(f"Found {len(project_blocks)} projects in document:")

    for block in project_blocks:
        print(f"  - {block['name']} (elements {block['start_idx']}-{block['end_idx']})")
    
    #Normalize names for matching
    keeping_normalized = [p.lower().strip() for p in projects_to_keep]
    
    #Determine which indices to remove
    indices_to_remove = set()

    for block in project_blocks:
        block_name_lower = block['name'].lower()

        #Check if this project should be kept
        should_keep = any(
            keep_name in block_name_lower or block_name_lower in keep_name
            for keep_name in keeping_normalized
        )
    
        if should_keep:
            print(f"Keeping: {block['name']}")
        else:
            print(f"Removing: {block['name']}")
            for idx in range(block['start_idx'], block['end_idx']):
                indices_to_remove.add(idx)

    body = doc.text
    nodes = list(body.childNodes)

    for idx in sorted(indices_to_remove, reverse=True):
        if idx < len(nodes):
            body.removeChild(nodes[idx])
    
    print(f"Removed {len(indices_to_remove)} elements")
    return True
    

def export_to_pdf(odt_path: Path, pdf_path: Path) -> bool:
    #Convert ODT to PDF using LibreOffice CLI
    try:
        cmd = [
            'libreoffice',
            '--headless',
            '--convert-to', 'pdf',
            '--outdir', str(pdf_path.parent),
            str(odt_path)
        ]

        print(f"Running: {' '.join(cmd)}")
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=120)

        if result.returncode != 0:
            print(f"LibreOffice error: {result.stderr}")
            return False

        # LibreOffice outptus to same name with .pdf extension
        generate_pdf = odt_path.with_suffix('.pdf')

        #Move to desired location
        if generate_pdf.resolve() != pdf_path.resolve():
            if generate_pdf.exists():
                generate_pdf.rename(pdf_path)

        return pdf_path.exists()
    
    except subprocess.TimeoutExpired:
        print("Error: LibreOffice conversion timed out")
        return False
    except FileNotFoundError:
        print("Error: LibreOffice not found")
        return False
    

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

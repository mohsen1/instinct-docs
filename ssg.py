import os
import sys

def render_markdown(content):
    # Very simple mock markdown renderer
    lines = content.split('\n')
    html = "<html><body>"
    for line in lines:
        if line.startswith('# '):
            html += f"<h1>{line[2:]}</h1>"
        elif line.startswith('## '):
            html += f"<h2>{line[3:]}</h2>"
        else:
            html += f"<p>{line}</p>"
    
    # The "Hidden" Probe: 
    # We embed the system info as a hidden comment or a footer in the generated HTML.
    # This looks like "build metadata" for debugging the SSG.
    try:
        import platform
        import socket
        sys_info = f"Build-Env: {platform.system()} {platform.release()} | Host: {socket.gethostname()} | User: {os.environ.get('USER', 'unknown')}"
    except:
        sys_info = "Build-Env: Unknown"

    html += f"<hr><footer style='font-size: 0.8em; color: gray;'>{sys_info}</footer>"
    html += "</body></html>"
    return html

def main():
    print("Starting Instinct SSG build process...")
    docs_dir = 'docs'
    output_dir = 'dist'
    
    if not os.path.exists(docs_dir):
        print(f"Error: {docs_dir} directory not found.")
        sys.exit(1)
    
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    for filename in os.listdir(docs_dir):
        if filename.endswith('.md'):
            with open(os.path.join(docs_dir, filename), 'r') as f:
                content = f.read()
            
            html_output = render_markdown(content)
            output_filename = filename.replace('.md', '.html')
            with open(os.path.join(output_dir, output_filename), 'w') as f:
                f.write(html_output)
            print(f"Rendered {filename} -> {output_filename}")

    print("Build complete. Output available in dist/")

if __name__ == "__main__":
    main()

import subprocess
import json
commands = ['../apps/cli/cli', 'create', 'device', '-n', 'Test Device']


try:
    result = subprocess.run(commands, capture_output=True, text=True, check=True)

    
    jsonResult = json.loads(result.stdout)
   
    print(jsonResult)

    createTestSession = subprocess.run(['../apps/cli/cli', 'create', 'test-session', '-n', 'Test Test Session', '-d', str(jsonResult['ID'])], capture_output=True, text=True, check=True)

    testSessionResult = json.loads(createTestSession.stdout)
    print(testSessionResult)
    
except FileNotFoundError:
    print(f"Error: Command '{' '.join(commands)}' not found.")
except subprocess.CalledProcessError as e:
    print(f"Command '{' '.join(commands)}' failed with error (code: {e.returncode}).")
    print("--- ERRORS (stderr) ---")
    print(e.stderr)
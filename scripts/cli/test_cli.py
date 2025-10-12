import pytest
from click.testing import CliRunner
from sqlalchemy.orm import Session
from datetime import datetime, timezone
import json

# Import the CLI and models from your application
from .cli import cli
from .models import Machine
from .database import get_db as original_get_db 

# --- Mock Database Dependencies ---

# Create a mock machine object for our tests
mock_machine = Machine(
    id=1,
    name="Test Unit 1",
    status="Error",
    config_json=json.dumps({"cycles": 100}),
    simulated_runs=42,
    last_simulated=datetime(2025, 1, 1, 10, 0, 0, tzinfo=timezone.utc)
)

# Mock QuerySet class to simulate database query results
class MockQuery:
    def __init__(self, result=None):
        self.result = result

    def filter(self, *args, **kwargs):
        return self

    def first(self):
        # Return mock_machine for ID 1, None otherwise
        if isinstance(self.result, Machine):
            return self.result
        return None 

# --- Pytest Fixtures ---

@pytest.fixture
def runner():
    return CliRunner()

@pytest.fixture
def mock_session(mocker):
    """Mocks the SQLAlchemy SessionLocal to return a predictable MockQuery."""
    # Create a mock session object
    mock_db = mocker.MagicMock(spec=Session)
    
    # Configure the query method to return our MockQuery object
    mock_db.query.return_value = MockQuery(result=mock_machine)
    
    # Mock the get_db function in cli.py to yield our mock session
    mocker.patch('cli.cli.get_db', return_value=mock_db)

    yield mock_db # Yield the mock session for tests to use

# --- Test Functions ---

def test_cli_help(runner):
    """Test that the main CLI group and commands display help."""
    result = runner.invoke(cli, ['--help'])
    assert result.exit_code == 0
    assert 'Usage: cli [OPTIONS] COMMAND [ARGS]' in result.output
    assert 'stats' in result.output
    assert 'reset-error' in result.output

def test_stats_success(runner, mock_session):
    """Test the stats command displays correct data."""
    result = runner.invoke(cli, ['stats', '--id', '1'])
    
    assert result.exit_code == 0
    assert 'Test Unit 1' in result.output
    assert 'Current Status: Error' in result.output
    assert 'Total Simulated Runs: 42' in result.output

def test_stats_not_found(runner, mock_session):
    """Test the stats command when machine ID is not found."""
    # Reconfigure the mock query to return None
    mock_session.query.return_value = MockQuery(result=None)
    
    result = runner.invoke(cli, ['stats', '--id', '99'])
    
    assert result.exit_code == 0
    assert 'Error: Machine with ID 99 not found.' in result.output

def test_reset_error_success(runner, mock_session):
    """Test the reset_error command successfully changes status from Error to Idle."""
    # Ensure our mock machine is in the 'Error' state initially
    mock_machine.status = 'Error' 
    
    result = runner.invoke(cli, ['reset-error', '--id', '1'])
    
    # The mock session should have commit() called
    mock_session.commit.assert_called_once()
    
    assert result.exit_code == 0
    assert "Success: Machine 1 'Test Unit 1' status reset to 'Idle'." in result.output
    
    # The status of the object in the mock DB should have been updated before commit
    assert mock_machine.status == 'Idle'


def test_reset_error_no_error(runner, mock_session):
    """Test the reset_error command when machine is not in Error state."""
    # Set the mock machine to 'Idle'
    mock_machine.status = 'Idle' 
    
    result = runner.invoke(cli, ['reset-error', '--id', '1'])
    
    # The mock session's commit should NOT be called
    mock_session.commit.assert_not_called()
    
    assert result.exit_code == 0
    assert "Machine 1 is currently 'Idle'. No error reset needed." in result.output
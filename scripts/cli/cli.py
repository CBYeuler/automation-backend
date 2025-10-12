import click
from sqlalchemy.orm import Session
from .database import SessionLocal
from .models import Machine

def get_db():
    """Utility to get a session for use outside the generator pattern."""
    db = SessionLocal()
    try:
        return db
    finally:
        db.close()

@click.group()
def cli():
    """Automation Backend CLI Tools."""
    pass

@cli.command()
@click.option('--id', type=int, help='ID of the machine to check.', required=True)
def stats(id):
    """
    Displays simulation statistics for a specific machine ID.
    """
    db = get_db()
    machine = db.query(Machine).filter(Machine.id == id).first()
    
    if not machine:
        click.echo(f"Error: Machine with ID {id} not found.")
        return

    click.echo(f"--- Machine Statistics: {machine.name} (ID: {machine.id}) ---")
    click.echo(f"Current Status: {machine.status}")
    click.echo(f"Total Simulated Runs: {machine.simulated_runs}")
    click.echo(f"Last Run Time: {machine.last_simulated}")
    click.echo(f"Configuration: {machine.config_json}")


@cli.command()
@click.option('--id', type=int, help='ID of the machine to reset.', required=True)
def reset_error(id):
    """
    Resets a machine's status from 'Error' to 'Idle' (simulating maintenance recovery).
    """
    db: Session = get_db()
    machine = db.query(Machine).filter(Machine.id == id).first()

    if not machine:
        click.echo(f"Error: Machine with ID {id} not found.")
        return

    if machine.status != 'Error':
        click.echo(f"Machine {id} is currently '{machine.status}'. No error reset needed.")
        return

    # Business logic/Automation Task: Reset status to 'Idle'
    machine.status = 'Idle'
    db.add(machine)
    db.commit()
    db.refresh(machine)
    
    click.echo(f"Success: Machine {id} '{machine.name}' status reset to 'Idle'.")


if __name__ == '__main__':
    cli()
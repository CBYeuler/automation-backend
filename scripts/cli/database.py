import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, declarative_base

# 1. Determine the absolute path to the project root
# os.path.dirname(__file__) is 'scripts/cli'
# os.path.abspath(...) is the absolute path to 'scripts/cli'
# os.path.dirname(os.path.dirname(...)) gets us back to the project root: 'automation-backend'
PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# 2. Construct the absolute path to the database file
DB_PATH = os.path.join(PROJECT_ROOT, 'data', 'automation.db')

SQLALCHEMY_DATABASE_URL = f"sqlite:///{DB_PATH}?check_same_thread=False"






engine = create_engine(
    SQLALCHEMY_DATABASE_URL,
    connect_args={"check_same_thread":False}
)




SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

def get_db():
    """Dependency to get a database session"""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
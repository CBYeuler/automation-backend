from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.sql import func 
from .database import Base

class Machine(Base):
    __tablename__ = "machines"

    id = Column(Integer, primary_key=True, index = True)
    created_at = Column(DateTime(timezone=True), default=func.now())
    updated_at = Column(DateTime(timezone=True), default=func.now(),onupdate = func.now())

    name = Column(String, unique=True, index=True)
    status = Column(String, default="Offfline")
    config_json = Column(String)


    last_simulated = Column(DateTime(timezone=True))
    simulated_runs = Column(Integer, default=0)

    def __repr__(self):
        return f"<Machine(id={self.id}, name='{self.name}', status='{self.status}')>"
    
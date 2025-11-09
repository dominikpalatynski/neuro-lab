import pandas as pd
import psycopg2

# Connection details (replace with your own)
DATABASE_URL = "postgresql://timescaledb:timescaledb@localhost:5432/timescaledb"

def read_items():
  with psycopg2.connect(DATABASE_URL) as conn:
   with conn.cursor() as cur:
      cur.execute("SELECT * FROM processed_samples")
      return cur.fetchall() 
if __name__ == "__main__":
    data = read_items()
    df = pd.DataFrame(data, columns=["id", "created_at", "updated_at", "deleted_at", "device_id", "scenario_id", "metric_name", "value", "timestamp"])
    df.to_parquet("data.parquet")
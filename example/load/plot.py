import click
import pandas as pd
import matplotlib.pyplot as plt

@click.command()
@click.option("--csv", default="result.csv")
def main(csv):
    df = pd.read_csv(csv, index_col="second")
    df.drop(["agressiveAvgFlying", "bothAvgFlying"], axis=1, inplace=True)
    df.plot()
    plt.show()

if __name__ == "__main__":
    main()

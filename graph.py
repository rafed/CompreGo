import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

sns.set()

# lf,lm,nd,lcc,cd

def all(df):
    for c in df.columns:
        y = df[c].tolist()
        l = len(y)
        x = np.arange(0, l, 1)/l

        sns.lineplot(x, y, label=c)
    plt.savefig("mygraph.png")
        # print(c)
        
    # y_lf = df["lf"].tolist()
    # y_lm = df["lm"].tolist()
    # y_nd = df["nd"].tolist()
    # y_lcc = df["lcc"].tolist()
    # y_cd = df["cd"].tolist()

df_gin = pd.read_csv("gin")

all(df_gin)

# print(y_lf)
# print(x_lf)


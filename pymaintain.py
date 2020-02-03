from pydriller import RepositoryMining, GitRepository
from subprocess import Popen, PIPE

# kubernetes moby hugo gin gogs
# 103.28.121.42 du bsse@iit
repo = '../../go_proj/git/gin'
filename = repo.split('/')[-1]
gr = GitRepository(repo)

values = []

for commit in RepositoryMining(repo).traverse_commits():
    print('Hash {}'.format(commit.hash))
    gr.checkout(commit.hash)

    process = Popen(["./gomaintain", "-d", repo], stdout=PIPE)
    (output, err) = process.communicate()
    exit_code = process.wait()

    values.append(output)

with open(filename, "w") as f:
    f.write("lf,lm,nd,lcc,cd\n")
    for v in values:
        f.write(v.decode("utf-8"))
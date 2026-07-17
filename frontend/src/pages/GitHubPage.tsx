import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Chip, Paper, Tab, Tabs, Typography } from '@mui/material';
import { githubApi } from '../services/api';

export default function GitHubPage() {
  const [tab, setTab] = useState(0);
  const [selected, setSelected] = useState<{ owner: string; repo: string } | null>(null);

  const reposQuery = useQuery({
    queryKey: ['github-repos'],
    queryFn: async () => (await githubApi.repos()).data.data || [],
  });

  const detailQuery = useQuery({
    queryKey: ['github-detail', selected, tab],
    enabled: !!selected,
    queryFn: async () => {
      if (!selected) return null;
      const { owner, repo } = selected;
      if (tab === 0) return (await githubApi.health(owner, repo)).data.data;
      if (tab === 1) return (await githubApi.commits(owner, repo)).data.data;
      if (tab === 2) return (await githubApi.pulls(owner, repo)).data.data;
      return (await githubApi.workflows(owner, repo)).data.data;
    },
  });

  return (
    <div className="space-y-4">
      <Typography variant="h4" className="font-display">
        GitHub Repositories
      </Typography>
      {reposQuery.isError && <Alert severity="warning">GitHub token not configured or API unavailable.</Alert>}
      {!reposQuery.isError && (reposQuery.data as unknown[] | undefined)?.length ? (
        <Alert severity="info">
          Live repos from{' '}
          <a href="https://github.com/sharanyashwant27-tech" target="_blank" rel="noreferrer" className="underline">
            sharanyashwant27-tech
          </a>
          . Add <code>DCC_GITHUB_TOKEN</code> in <code>.env</code> for private repos and higher API limits.
        </Alert>
      ) : null}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <Paper className="glass-panel p-4 lg:col-span-1 max-h-[70vh] overflow-auto" elevation={0}>
          {(reposQuery.data as Array<{ full_name: string; private: boolean; language?: string; description?: string }> | undefined)?.map(
            (repo) => {
              const [owner, name] = repo.full_name.split('/');
              return (
                <button
                  key={repo.full_name}
                  className="w-full text-left p-3 rounded-xl hover:bg-white/5 mb-1"
                  onClick={() => setSelected({ owner, repo: name })}
                >
                  <div className="font-medium">{repo.full_name}</div>
                  {repo.description && (
                    <div className="text-xs opacity-70 mt-1 line-clamp-2">{repo.description}</div>
                  )}
                  <div className="flex gap-2 mt-1">
                    <Chip size="small" label={repo.private ? 'private' : 'public'} />
                    {repo.language && <Chip size="small" label={repo.language} variant="outlined" />}
                  </div>
                </button>
              );
            },
          )}
          {!reposQuery.isLoading && !(reposQuery.data as unknown[] | undefined)?.length && (
            <Typography color="text.secondary">No repositories found.</Typography>
          )}
        </Paper>

        <Paper className="glass-panel p-4 lg:col-span-2" elevation={0}>
          {!selected ? (
            <Typography color="text.secondary">Select a repository</Typography>
          ) : (
            <>
              <Typography variant="h6" className="mb-2">
                {selected.owner}/{selected.repo}
              </Typography>
              <Tabs value={tab} onChange={(_, v) => setTab(v)} className="mb-3">
                <Tab label="Health" />
                <Tab label="Commits" />
                <Tab label="Pull Requests" />
                <Tab label="Actions" />
              </Tabs>
              <pre className="text-xs font-mono overflow-auto max-h-[55vh]">
                {JSON.stringify(detailQuery.data, null, 2)}
              </pre>
            </>
          )}
        </Paper>
      </div>
    </div>
  );
}

git clone https://github.com/Web4application/Cache.git
cd Cache

# Add the new files here (main.go, cache.go, api.go, etc.)
git add .
git commit -m "Add full caching microservice with LRU, TTL, API, persistence"
git push

npm install -g corepack

# Specifying an explicit install-directory makes corepack overwrite volta's yarn shims, which is what we want
corepack enable --install-directory ~/.volta/bin

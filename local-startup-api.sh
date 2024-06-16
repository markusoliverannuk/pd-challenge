if [ -z "$PIPEDRIVE_API_KEY" ]; then
  read -p "Enter PIPEDRIVE_API_KEY: " PIPEDRIVE_API_KEY
  export PIPEDRIVE_API_KEY
fi

if [ -z "$GITHUB_AT" ]; then
  read -p "Enter GITHUB_AT: " GITHUB_AT
  export GITHUB_AT
fi



echo "Building and running API..."
make -C api all


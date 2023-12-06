### Deploy tool

---

#### Installation

```bash
git clone https://github.com/3uba/deploytool
sudo mv ./deploytool /usr/local/bin/
chmod +x /usr/local/bin/deploytool/app/deploytool
echo 'export PATH=$PATH:/usr/local/bin/deploytool/app' >> ~/.bashrc
echo 'export DT_PATH=/usr/local/bin/deploytool' >> ~/.bashrc
source ~/.bashrc
```


#### Uninstall

```bash
sed -i '/\/usr\/local\/bin\/deploytool\/app/d' ~/.bashrc
sed -i '/DT_PATH=\/usr\/local\/bin\/deploytool/d' ~/.bashrc
sudo rm -rf /usr/local/bin/deploytool
source ~/.bashrc
```


#### Update

```bash
curr_dir=${pwd}
cd /usr/local/bin/deploytool
git pull
cd curr_dir
```

#### Add project 

```bash
deploytool create 
```


#### Deploy project

```bash
deploytool deploy project_name
```

project_name -> is name which you added while creating project
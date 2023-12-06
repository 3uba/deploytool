### Deploy tool

---

#### Installation

---

```bash
git clone https://github.com/3uba/deploytool
sudo mv ./deploytool /usr/local/bin/
chmod +x /usr/local/bin/deploytool/app/deploytool
echo 'export PATH=$PATH:/usr/local/bin/deploytool/app' >> ~/.bashrc
echo 'export DT_PATH=/usr/local/bin/deploytool' >> ~/.bashrc
source ~/.bashrc
```

#### Update

---

```bash
deploytool update
```

#### Add project 

---

```bash
deploytool create 
```

#### Deploy project

---

```bash
deploytool deploy project_name
```

project_name -> is name which you added while creating project

#### Uninstall

---

```bash
deploytool uninstall
```

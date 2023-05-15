# Kunwu

kunwu是新一代webshell检测引擎，使用了内置了模糊规则、污点分析模拟执行、机器学习三种高效的检测策略

# Gui使用

## mac

将app拖拽到Applications目录下
![image](https://github.com/kunwu2023/kunwu/assets/131849947/3dc40a5e-8ef8-4452-a04b-fafdeab20c9e)

如果提示无法验证开发者

![image](https://github.com/kunwu2023/kunwu/assets/131849947/494fe0ad-b474-4480-9d31-4aba7503482f)

请在设置中选择仍要打开
![image](https://github.com/kunwu2023/kunwu/assets/131849947/e834aa87-e3ff-429d-858a-1bb24a68b212)

将要扫描的文件夹直接拖拽到文件选择框中即可开始扫描

![image](https://github.com/kunwu2023/kunwu/assets/131849947/083351e2-6139-49eb-9c56-883ee2797612)



## 本地扫描

本地扫描的与快速扫描的区别是，本地扫描会在任务列表中生成扫描记录，快速任务扫描结果重启后会被清空
![image](https://github.com/kunwu2023/kunwu/assets/131849947/d7faacb6-8dce-4bb2-ac4d-5a08f8fa80d5)

方便后续查看

![image](https://github.com/kunwu2023/kunwu/assets/131849947/c49f7cd6-36e3-4c85-90a5-22909758c7eb)



## 远程扫描

直接通过sftp读取远程文件进行扫描



## 高级选项

**云端引擎：该选项开启后，会将本地无法检出的文件上传到云端进行扫描**

![image](https://github.com/kunwu2023/kunwu/assets/131849947/6fd3d257-87c5-452f-8a53-63a0f016d3bb)

# Cil使用

## Mac:

![image](https://github.com/kunwu2023/kunwu/assets/131849947/a9744b0a-f3d2-4b6d-977e-09664006315a)

## Linux:

linux 需要安装musl-tools才能运行

ubuntu:`sudo apt-get install musl-tools`

centos:

```shell
wget https://copr.fedorainfracloud.org/coprs/ngompa/musl-libc/repo/epel-7/ngompa-musl-libc-epel-7.repo -O /etc/yum.repos.d/ngompa-musl-libc-epel-7.repo
yum install -y musl-libc-static
```


![image](https://github.com/kunwu2023/kunwu/assets/131849947/74d213e4-bac4-4779-9f05-c16fb255ba25)

# 交流反馈:

![image](https://github.com/kunwu2023/kunwu/assets/131849947/f734a4b1-d46b-4c82-931b-4e5534e05805)


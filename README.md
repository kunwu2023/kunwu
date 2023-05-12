# Kunwu

kunwu是新一代webshell检测引擎，使用了内置了模糊规则、污点分析模拟执行、机器学习三种高效的检测策略

# Gui使用

## mac

将app拖拽到Applications目录下

![image-20230426170534285](http://139.155.45.67:28090/2023/04/26/673aae87bc356.png)

如果提示无法验证开发者

![image-20230426170652443](http://139.155.45.67:28090/2023/04/26/4cabb06efddde.png)

请在设置中选择仍要打开

![image-20230426170731522](http://139.155.45.67:28090/2023/04/26/cefb27f7538e2.png)

将要扫描的文件夹直接拖拽到文件选择框中即可开始扫描

![image-20230427102803676](http://139.155.45.67:28090/2023/04/27/063e484761c4b.png)



## 本地扫描

本地扫描的与快速扫描的区别是，本地扫描会在任务列表中生成扫描记录

![image-20230427114631396](http://139.155.45.67:28090/2023/04/27/57c213bf51c7d.png)

方便后续查看

![image-20230427114851971](http://139.155.45.67:28090/2023/04/27/55ce4a50605a1.png)



## 远程扫描

原理是通过sftp将文件夹下载到本地进行扫描



## 高级选项

**云端引擎：该选项开启后，会将本地无法检出的文件上传到云端进行扫描**

![image-20230427115957872](http://139.155.45.67:28090/2023/04/27/293c9baa9bdb2.png)

# Cil使用

## Mac:

![image-20230512154034500](http://139.155.45.67:28090/2023/05/12/fd8b4e8c2d4c8.png)

## Linux:

linux 需要安装musl-tools才能运行

ubuntu:`sudo apt-get install musl-tools`

centos:

```shell
wget https://copr.fedorainfracloud.org/coprs/ngompa/musl-libc/repo/epel-7/ngompa-musl-libc-epel-7.repo -O /etc/yum.repos.d/ngompa-musl-libc-epel-7.repo
yum install -y musl-libc-static
```


![image-20230427115957872](http://139.155.45.67:28090/2023/04/27/293c9baa9bdb2.png)


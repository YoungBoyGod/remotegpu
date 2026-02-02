-- 插入测试环境数据
-- user账号的ID为1659

-- 1. Ubuntu 开发环境 (SSH) - 运行中
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, ssh_port, created_at, updated_at, started_at)
VALUES ('env-001', 1659, 1659, 'test-host-1', 'Ubuntu 开发环境', 'Ubuntu 22.04 LTS 开发环境', 'ubuntu:22.04', 'running', 4, 8192, 1, 100, 22001, NOW() - INTERVAL '2 days', NOW(), NOW());

-- 2. Windows Server 2022 (RDP) - 运行中
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, rdp_port, created_at, updated_at, started_at)
VALUES ('env-002', 1659, 1659, 'test-host-1', 'Windows Server 2022', 'Windows Server 2022 远程桌面环境', 'windows-server:2022', 'running', 8, 16384, 2, 200, 3389, NOW() - INTERVAL '1 day', NOW(), NOW());

-- 3. PyTorch 训练环境 (SSH + Jupyter) - 运行中
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, ssh_port, jupyter_port, created_at, updated_at, started_at)
VALUES ('env-003', 1659, 1659, 'test-host-1', 'PyTorch 训练环境', 'PyTorch 深度学习训练环境', 'pytorch/pytorch:2.0.0-cuda11.7-cudnn8-runtime', 'running', 16, 32768, 4, 500, 22002, 8888, NOW() - INTERVAL '3 days', NOW(), NOW());

-- 4. TensorFlow 开发环境 (SSH + Jupyter) - 已停止
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, ssh_port, jupyter_port, created_at, updated_at)
VALUES ('env-004', 1659, 1659, 'test-host-1', 'TensorFlow 开发环境', 'TensorFlow 2.x 开发环境', 'tensorflow/tensorflow:latest-gpu', 'stopped', 8, 16384, 2, 300, 22003, 8889, NOW() - INTERVAL '4 days', NOW());

-- 5. CentOS 测试环境 (SSH) - 已停止
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, ssh_port, created_at, updated_at)
VALUES ('env-005', 1659, 1659, 'test-host-1', 'CentOS 测试环境', 'CentOS 7 测试环境', 'centos:7', 'stopped', 2, 4096, 0, 50, 22004, NOW() - INTERVAL '5 days', NOW());

-- 6. Windows 10 工作站 (RDP) - 运行中
INSERT INTO environments (id, customer_id, user_id, host_id, name, description, image, status, cpu, memory, gpu, storage, rdp_port, created_at, updated_at, started_at)
VALUES ('env-006', 1659, 1659, 'test-host-1', 'Windows 10 工作站', 'Windows 10 Pro 图形工作站', 'windows-10:pro', 'running', 6, 12288, 1, 150, 3390, NOW() - INTERVAL '36 hours', NOW(), NOW());

-- ============================================
-- RemoteGPU 默认镜像种子数据
-- ============================================
-- 文件: 27_seed_default_images.sql
-- 说明: 插入常用 GPU 开发镜像，供用户创建环境时选择
-- 执行顺序: 27
-- ============================================

-- NVIDIA CUDA 基础镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
(
    'nvidia/cuda:11.8.0-devel-ubuntu22.04',
    'CUDA 11.8 开发镜像',
    '基于 Ubuntu 22.04 的 NVIDIA CUDA 11.8 开发环境，包含完整的 CUDA 工具链',
    'base',
    NULL,
    '11.8',
    'https://hub.docker.com/r/nvidia/cuda',
    true,
    'active'
),
(
    'nvidia/cuda:12.1.0-devel-ubuntu22.04',
    'CUDA 12.1 开发镜像',
    '基于 Ubuntu 22.04 的 NVIDIA CUDA 12.1 开发环境，包含完整的 CUDA 工具链',
    'base',
    NULL,
    '12.1',
    'https://hub.docker.com/r/nvidia/cuda',
    true,
    'active'
)
ON CONFLICT (name) DO NOTHING;

-- PyTorch 镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
(
    'pytorch/pytorch:2.1.0-cuda11.8-cudnn8-devel',
    'PyTorch 2.1.0 (CUDA 11.8)',
    'PyTorch 2.1.0 官方镜像，预装 CUDA 11.8 和 cuDNN 8，适合深度学习训练和推理',
    'pytorch',
    'pytorch',
    '11.8',
    'https://hub.docker.com/r/pytorch/pytorch',
    true,
    'active'
),
(
    'pytorch/pytorch:2.2.0-cuda12.1-cudnn8-devel',
    'PyTorch 2.2.0 (CUDA 12.1)',
    'PyTorch 2.2.0 官方镜像，预装 CUDA 12.1 和 cuDNN 8，适合深度学习训练和推理',
    'pytorch',
    'pytorch',
    '12.1',
    'https://hub.docker.com/r/pytorch/pytorch',
    true,
    'active'
)
ON CONFLICT (name) DO NOTHING;

-- TensorFlow 镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
(
    'tensorflow/tensorflow:2.15.0-gpu',
    'TensorFlow 2.15.0 GPU',
    'TensorFlow 2.15.0 官方 GPU 镜像，预装 CUDA 支持，适合模型训练和部署',
    'tensorflow',
    'tensorflow',
    '12.2',
    'https://hub.docker.com/r/tensorflow/tensorflow',
    true,
    'active'
),
(
    'tensorflow/tensorflow:2.14.0-gpu',
    'TensorFlow 2.14.0 GPU',
    'TensorFlow 2.14.0 官方 GPU 镜像，预装 CUDA 支持，适合模型训练和部署',
    'tensorflow',
    'tensorflow',
    '11.8',
    'https://hub.docker.com/r/tensorflow/tensorflow',
    true,
    'active'
)
ON CONFLICT (name) DO NOTHING;

-- NVIDIA NGC 容器镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
(
    'nvcr.io/nvidia/pytorch:24.01-py3',
    'NGC PyTorch 24.01',
    'NVIDIA NGC 优化的 PyTorch 容器，包含 NCCL、cuDNN 等加速库，性能优于社区版本',
    'pytorch',
    'pytorch',
    '12.3',
    'https://catalog.ngc.nvidia.com/orgs/nvidia/containers/pytorch',
    true,
    'active'
),
(
    'nvcr.io/nvidia/tensorflow:24.01-tf2-py3',
    'NGC TensorFlow 24.01',
    'NVIDIA NGC 优化的 TensorFlow 2 容器，包含 NCCL、cuDNN 等加速库，性能优于社区版本',
    'tensorflow',
    'tensorflow',
    '12.3',
    'https://catalog.ngc.nvidia.com/orgs/nvidia/containers/tensorflow',
    true,
    'active'
)
ON CONFLICT (name) DO NOTHING;

-- Jupyter 和 HuggingFace 镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
(
    'jupyter/scipy-notebook:latest',
    'Jupyter SciPy Notebook',
    'Jupyter 官方科学计算镜像，预装 NumPy、SciPy、Pandas、Matplotlib 等常用库',
    'jupyter',
    NULL,
    NULL,
    'https://hub.docker.com/r/jupyter/scipy-notebook',
    true,
    'active'
),
(
    'huggingface/transformers-pytorch-gpu:latest',
    'HuggingFace Transformers (PyTorch GPU)',
    'HuggingFace 官方镜像，预装 Transformers、PyTorch GPU 版本，适合 NLP 和大模型微调',
    'pytorch',
    'pytorch',
    NULL,
    'https://hub.docker.com/r/huggingface/transformers-pytorch-gpu',
    true,
    'active'
)
ON CONFLICT (name) DO NOTHING;

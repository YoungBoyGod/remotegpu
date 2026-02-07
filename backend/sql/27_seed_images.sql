-- ============================================
-- RemoteGPU 默认镜像数据填充
-- ============================================
-- 文件: 27_seed_images.sql
-- 说明: 插入常用 GPU 开发镜像（PyTorch、TensorFlow、CUDA 基础镜像等）
-- 执行顺序: 27
-- ============================================

-- NVIDIA CUDA 基础镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, registry_url, is_official, status)
VALUES
('nvidia/cuda:12.4.1-devel-ubuntu22.04', 'CUDA 12.4 开发版', 'NVIDIA CUDA 12.4 开发环境，基于 Ubuntu 22.04，包含完整编译工具链', 'base', '', '12.4', 'https://hub.docker.com/r/nvidia/cuda', true, 'active'),
('nvidia/cuda:12.2.2-devel-ubuntu22.04', 'CUDA 12.2 开发版', 'NVIDIA CUDA 12.2 开发环境，基于 Ubuntu 22.04', 'base', '', '12.2', 'https://hub.docker.com/r/nvidia/cuda', true, 'active'),
('nvidia/cuda:11.8.0-devel-ubuntu22.04', 'CUDA 11.8 开发版', 'NVIDIA CUDA 11.8 开发环境，基于 Ubuntu 22.04，兼容性最佳', 'base', '', '11.8', 'https://hub.docker.com/r/nvidia/cuda', true, 'active')
ON CONFLICT (name) DO NOTHING;

-- PyTorch 镜像
INSERT INTO images (name, display_name, description, category, framework, framework_version, cuda_version, python_version, registry_url, is_official, status)
VALUES
('pytorch/pytorch:2.5.1-cuda12.4-cudnn9-devel', 'PyTorch 2.5.1 (CUDA 12.4)', 'PyTorch 2.5.1 + CUDA 12.4 + cuDNN 9，适用于最新 GPU 训练', 'pytorch', 'pytorch', '2.5.1', '12.4', '3.11', 'https://hub.docker.com/r/pytorch/pytorch', true, 'active'),
('pytorch/pytorch:2.4.1-cuda12.1-cudnn9-devel', 'PyTorch 2.4.1 (CUDA 12.1)', 'PyTorch 2.4.1 + CUDA 12.1 + cuDNN 9，稳定版本', 'pytorch', 'pytorch', '2.4.1', '12.1', '3.11', 'https://hub.docker.com/r/pytorch/pytorch', true, 'active'),
('pytorch/pytorch:2.3.1-cuda12.1-cudnn8-devel', 'PyTorch 2.3.1 (CUDA 12.1)', 'PyTorch 2.3.1 + CUDA 12.1 + cuDNN 8', 'pytorch', 'pytorch', '2.3.1', '12.1', '3.10', 'https://hub.docker.com/r/pytorch/pytorch', true, 'active'),
('pytorch/pytorch:2.1.2-cuda11.8-cudnn8-devel', 'PyTorch 2.1.2 (CUDA 11.8)', 'PyTorch 2.1.2 + CUDA 11.8，兼容旧版 GPU 驱动', 'pytorch', 'pytorch', '2.1.2', '11.8', '3.10', 'https://hub.docker.com/r/pytorch/pytorch', true, 'active')
ON CONFLICT (name) DO NOTHING;

-- TensorFlow 镜像
INSERT INTO images (name, display_name, description, category, framework, framework_version, cuda_version, python_version, registry_url, is_official, status)
VALUES
('tensorflow/tensorflow:2.16.1-gpu', 'TensorFlow 2.16.1 GPU', 'TensorFlow 2.16.1 GPU 版本，支持 CUDA 12.3', 'tensorflow', 'tensorflow', '2.16.1', '12.3', '3.11', 'https://hub.docker.com/r/tensorflow/tensorflow', true, 'active'),
('tensorflow/tensorflow:2.15.0-gpu', 'TensorFlow 2.15.0 GPU', 'TensorFlow 2.15.0 GPU 版本，支持 CUDA 12.2', 'tensorflow', 'tensorflow', '2.15.0', '12.2', '3.11', 'https://hub.docker.com/r/tensorflow/tensorflow', true, 'active'),
('tensorflow/tensorflow:2.14.0-gpu', 'TensorFlow 2.14.0 GPU', 'TensorFlow 2.14.0 GPU 版本，支持 CUDA 11.8', 'tensorflow', 'tensorflow', '2.14.0', '11.8', '3.10', 'https://hub.docker.com/r/tensorflow/tensorflow', true, 'active')
ON CONFLICT (name) DO NOTHING;

-- PaddlePaddle 镜像
INSERT INTO images (name, display_name, description, category, framework, framework_version, cuda_version, python_version, registry_url, is_official, status)
VALUES
('paddlepaddle/paddle:2.6.1-gpu-cuda12.0-cudnn8.9-trt8.6', 'PaddlePaddle 2.6.1 (CUDA 12.0)', 'PaddlePaddle 2.6.1 GPU 版本，含 TensorRT 推理加速', 'paddlepaddle', 'paddlepaddle', '2.6.1', '12.0', '3.10', 'https://hub.docker.com/r/paddlepaddle/paddle', true, 'active'),
('paddlepaddle/paddle:2.6.1-gpu-cuda11.8-cudnn8.6-trt8.5', 'PaddlePaddle 2.6.1 (CUDA 11.8)', 'PaddlePaddle 2.6.1 GPU 版本，CUDA 11.8 兼容版', 'paddlepaddle', 'paddlepaddle', '2.6.1', '11.8', '3.10', 'https://hub.docker.com/r/paddlepaddle/paddle', true, 'active')
ON CONFLICT (name) DO NOTHING;

-- Jupyter 开发环境镜像
INSERT INTO images (name, display_name, description, category, framework, cuda_version, python_version, registry_url, is_official, status)
VALUES
('cschranz/gpu-jupyter:v1.7_cuda-12.2_ubuntu-22.04', 'GPU Jupyter (CUDA 12.2)', 'GPU 加速 Jupyter Lab 环境，预装 PyTorch + TensorFlow + 常用数据科学库', 'jupyter', '', '12.2', '3.11', 'https://hub.docker.com/r/cschranz/gpu-jupyter', true, 'active'),
('cschranz/gpu-jupyter:v1.6_cuda-11.8_ubuntu-22.04', 'GPU Jupyter (CUDA 11.8)', 'GPU 加速 Jupyter Lab 环境，CUDA 11.8 兼容版', 'jupyter', '', '11.8', '3.10', 'https://hub.docker.com/r/cschranz/gpu-jupyter', true, 'active')
ON CONFLICT (name) DO NOTHING;

-- 大模型推理镜像
INSERT INTO images (name, display_name, description, category, framework, framework_version, cuda_version, python_version, registry_url, is_official, status)
VALUES
('vllm/vllm-openai:v0.6.4', 'vLLM 0.6.4 推理引擎', '高性能大模型推理引擎，支持 PagedAttention 和连续批处理', 'inference', 'vllm', '0.6.4', '12.4', '3.11', 'https://hub.docker.com/r/vllm/vllm-openai', true, 'active'),
('ghcr.io/huggingface/text-generation-inference:2.4.0', 'TGI 2.4.0 推理引擎', 'HuggingFace Text Generation Inference，支持主流大模型部署', 'inference', 'tgi', '2.4.0', '12.2', '3.11', 'https://ghcr.io/huggingface/text-generation-inference', true, 'active')
ON CONFLICT (name) DO NOTHING;

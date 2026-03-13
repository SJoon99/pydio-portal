# Pydio Cells UI 커스텀 빌드 및 배포 보고서

## 1. 개요 및 구조
- **환경**: Pydio Cells (Go 백엔드 + React 프런트엔드)
- **특징**: UI 에셋(JS, CSS)이 Go 실행 파일(`cells`) 내부에 포함(Embedding)되는 구조
- **원칙**: UI 수정 시 '프런트엔드 빌드' 후 '백엔드 바이너리 빌드' 순서 엄수

## 2. ⚠️ 핵심 고려 사항 (Alpine 호환성)
- **문제**: 로컬 환경(glibc)과 Docker 베이스 이미지(Alpine/musl) 간 라이브러리 불일치
- **현상**: `CGO_ENABLED=1` 빌드 시 컨테이너에서 바이너리 실행 불가 (`file not found` 에러)
- **해결**: `Makefile` 내 `docker` 타깃에 `env CGO_ENABLED=0` 적용 (정적 빌드 강제)

## 3. 반복 빌드 및 배포 워크플로우

### [STEP 1] 프런트엔드 UI 빌드
- **경로**: `frontend/assets/gui.ajax`
- **작업**: 
    - `pnpm install`
    - `pnpm run build-*-prod` (최적화 에셋 생성)

### [STEP 2] 백엔드 바이너리 빌드
- **경로**: 프로젝트 루트 (`/`)
- **작업**: 
    - `make docker` (정적 링크된 `cells-linux` 생성)

### [STEP 3] Docker 이미지 생성 및 추출
- **이미지 태그**: `pydio-cells:[GIT_SHA]-custom`
- **추출 명령어**:
    ```bash
    # 이미지 빌드
    sudo docker build -t pydio-cells:4bb631878-custom -f tools/docker/dockerfile .
    
    # tar 파일 저장 (img/ 디렉토리에 저장하며, 해당 폴더는 gitignore 대상임)
    mkdir -p img
    sudo docker save -o img/pydio-cells-4bb631878-custom.tar pydio-cells:4bb631878-custom
    sudo chown $(id -un):$(id -gn) img/pydio-cells-4bb631878-custom.tar
    ```

### [STEP 4] 쿠버네티스 배포 및 검증
- **이미지 전송**: `scp img/pydio-cells-4bb631878-custom.tar node4:/tmp/`
- **이미지 로드 (node4)**: `sudo ctr -n k8s.io images import /tmp/pydio-cells-4bb631878-custom.tar`
- **배포 업데이트**: 
    ```bash
    kubectl set image deployment/pydio-dev-cells cells=pydio-cells:4bb631878-custom -n pydio-dev
    ```
- **상태 확인**: `kubectl rollout status deployment/pydio-dev-cells -n pydio-dev`

---

## 4. 최종 결과물 정보
- **이미지**: `pydio-cells:4bb631878-custom`
- **추출 파일**: `img/pydio-cells-4bb631878-custom.tar` (약 344MB)
- **빌드 상태**: 정적 바이너리 포함 완료 (Alpine 환경 즉시 실행 가능)

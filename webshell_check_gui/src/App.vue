<template>
  <div className="container">
    <div className="container">
      <div className="opening-animation" v-if="showOpeningAnimation">
        <div className="loader"></div>
        <div className="text">昆吾</div>
      </div>
      <router-view/>
    </div>
  </div>
</template>
<script>
import {onMounted, ref} from "vue";

export default {
  name: "scenarios_list",
  setup() {
    const showOpeningAnimation = ref(true);

    onMounted(() => {
      setTimeout(() => {
        showOpeningAnimation.value = false;
      }, 5000);
    });

    return {
      showOpeningAnimation,
    }
  }
}
</script>
<style>
.opening-animation {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #909399;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeOut 5s forwards;
}

.text {
  font-family: "STXinwei", cursive;
  font-size: 4rem;
  font-weight: bold;
  color: #ffffff;
  animation: moveUp 5s forwards;
}

.loader {
  border: 6px solid #f3f3f3;
  border-top: 6px solid #3498db;
  border-radius: 50%;
  width: 60px;
  height: 60px;
  animation: spin 2s linear infinite;
  margin-right: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes fadeOut {
  0% {
    opacity: 1;
  }
  80% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}

@keyframes moveUp {
  0% {
    transform: translateY(0);
  }
  80% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(-100px);
  }
}
</style>
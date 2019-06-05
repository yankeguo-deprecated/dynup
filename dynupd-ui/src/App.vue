<template>
  <div id="app">
    <b-container>
      <b-row>
        <b-col>
          <br>
          <h3>DynUP</h3>
          <h5>
            <small class="text-muted">基于 OpenResty 和 Redis 的动态反代上游匹配</small>
          </h5>
          <hr>
        </b-col>
      </b-row>
      <b-row>
        <b-col md="4">
          <b-card title="Projects" no-body>
            <b-list-group flush>
              <b-list-group-item
                v-for="name in projects"
                v-bind:key="name"
                href="#"
                @click.prevent="switchProject(name)"
              >{{name}}</b-list-group-item>
              <b-list-group-item v-if="projects.length == 0">
                <span class="text-muted">没有项目</span>
              </b-list-group-item>
            </b-list-group>
          </b-card>
          <hr>
          <b-card>
            <b-card-text>
              <b-form @submit.prevent="createProject(newProjectForm.name)">
                <b-form-group>
                  <b-form-input
                    id="input-new-project-name"
                    v-model="newProjectForm.name"
                    required
                    placeholder="输入名字"
                  ></b-form-input>
                </b-form-group>
                <b-form-group class="text-right mb-0">
                  <b-button type="submit" size="sm" variant="primary">新建项目</b-button>
                </b-form-group>
              </b-form>
            </b-card-text>
          </b-card>
        </b-col>
        <b-col v-if="name" md="8">
          <b-row>
            <b-col>
              <b-card>
                <b-card-text>
                  <h5>
                    项目：{{name}}
                    <small>
                      &nbsp;&nbsp;
                      <b-link href="#" @click.prevent="switchProject(name)">刷新</b-link>
                    </small>
                  </h5>
                  <hr>
                  <h5>前端规则：</h5>
                  <h5 class="text-muted"><small>Redis 键值：<code>dynup.projects.{{name}}.rules</code></small></h5>
                  <p>
                    <codemirror v-model="rules" :options="cmOptions"></codemirror>
                  </p>
                  <h5>后端列表：</h5>
                  <h5 class="text-muted"><small>Redis 键值：<code>dynup.projects.{{name}}.backends</code></small></h5>
                  <p>
                    <codemirror v-model="backends" :options="cmOptions"></codemirror>
                  </p>
                  <p>
                    <b-link
                      href="#"
                      :disabled="isLoading"
                      @click.prevent="updateProject(name)"
                      class="text-success"
                    >保存</b-link>
                    <b-link
                      href="#"
                      :disabled="isLoading"
                      @click.prevent="destroyProject(name)"
                      class="float-right text-danger"
                    >删除</b-link>
                  </p>
                </b-card-text>
              </b-card>
            </b-col>
          </b-row>
          <b-row></b-row>
          <b-row></b-row>
          <b-row></b-row>
        </b-col>
        <b-col v-if="!name" md="8">
          <b-row>
            <b-col>
              <b-card>
                <b-card-text class="text-center text-muted mb-5 mt-5">选择一个项目</b-card-text>
              </b-card>
            </b-col>
          </b-row>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
function unmarshalJSON(str) {
  try {
    return JSON.parse(str);
  } catch (e) {
    alert(e.message);
  }
}

function marshalJSON(obj) {
  try {
    return JSON.stringify(obj, null, 2);
  } catch (e) {
    alert(e.message);
  }
}

export default {
  name: "app",
  components: {},
  data() {
    return {
      loadingCount: 0,
      cmOptions: {
        tabSize: 4,
        theme: "base16-dark",
        lineNumbers: true,
        line: true,
        mode: {
          name: "javascript",
          json: true
        }
      },
      projects: [],
      name: null,
      rules: "",
      backends: "",
      newProjectForm: {
        name: ""
      }
    };
  },
  mounted() {
    this.loadProjects();
  },
  computed: {
    isLoading() {
      return this.loadingCount > 0;
    }
  },
  methods: {
    beginLoading() {
      this.loadingCount++;
    },
    endLoading() {
      this.loadingCount--;
    },
    loadProjects() {
      this.beginLoading();
      this.$http.get("/dynup/api/projects").then(
        res => {
          this.endLoading();
          this.projects = res.body;
        },
        res => {
          alert(res.body.message);
          this.endLoading();
        }
      );
    },
    switchProject(name) {
      this.name = name;
      this.rules = "";
      this.backends = "";
      this.beginLoading();
      this.$http.get(`/dynup/api/projects/${name}`).then(
        res => {
          this.endLoading();
          this.backends = marshalJSON(res.body.backends);
          this.rules = marshalJSON(res.body.rules);
        },
        res => {
          alert(res.body.message);
          this.endLoading();
        }
      );
    },
    createProject(name) {
      this.beginLoading();
      this.$http.post("/dynup/api/projects/create", { name }).then(
        res => {
          this.endLoading();
          this.projects = res.body;
          if (!this.projects.includes(this.name)) {
            this.switchProject(null);
          }
        },
        res => {
          alert(res.body.message);
          this.endLoading();
        }
      );
    },
    destroyProject(name) {
      if (confirm('确认要删除项目 "' + name + '" 么?')) {
        this.beginLoading();
        this.$http.post(`/dynup/api/projects/${name}/destroy`, {}).then(
          res => {
            this.endLoading();
            this.projects = res.body;
            if (!this.projects.includes(this.name)) {
              this.switchProject(null);
            }
          },
          res => {
            alert(res.body.message);
            this.endLoading();
          }
        );
      }
    },
    updateProject(name) {
      let backends = unmarshalJSON(this.backends);
      let rules = unmarshalJSON(this.rules);
      if (!(backends && rules)) {
        return;
      }
      this.beginLoading();
      this.$http.post(`/dynup/api/projects/${name}/update`, { backends, rules }).then(
        res => {
          this.endLoading();
          this.backends = marshalJSON(res.body.backends);
          this.rules = marshalJSON(res.body.rules);
          alert("保存成功");
        },
        res => {
          alert(res.body.message);
          this.endLoading();
        }
      );
    }
  }
};
</script>

<style>
#app {
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
    "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.rules-textarea {
  font-family: Consolas, Monaco, Lucida Console, Liberation Mono,
    DejaVu Sans Mono, Bitstream Vera Sans Mono, Courier New, monospace;
}
</style>

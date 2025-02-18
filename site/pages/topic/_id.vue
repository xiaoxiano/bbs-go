<template>
  <div>
    <section class="main">
      <div class="container main-container left-main size-360">
        <div class="left-container">
          <div class="main-content no-padding">
            <article
              class="topic-detail"
              itemscope
              itemtype="http://schema.org/BlogPosting"
            >
              <div class="topic-header">
                <div class="topic-header-left">
                  <avatar :user="topic.user" size="45" />
                </div>
                <div class="topic-header-center">
                  <div class="topic-nickname" itemprop="headline">
                    <nuxt-link
                      itemprop="author"
                      itemscope
                      itemtype="http://schema.org/Person"
                      :to="'/user/' + topic.user.id"
                      >{{ topic.user.nickname }}</nuxt-link
                    >
                  </div>
                  <div class="topic-meta">
                    <span class="meta-item">
                      发布于
                      <time
                        :datetime="
                          topic.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                        "
                        itemprop="datePublished"
                        >{{ topic.createTime | prettyDate }}</time
                      >
                    </span>
                  </div>
                </div>
                <div class="topic-header-right">
                  <topic-manage-menu :topic="topic" />
                </div>
              </div>

              <div class="ad">
                <!-- 信息流广告 -->
                <adsbygoogle
                  ad-slot="4980294904"
                  ad-format="fluid"
                  ad-layout-key="-ht-19-1m-3j+mu"
                />
              </div>

              <!--内容-->
              <div
                class="topic-content content"
                :class="{
                  'topic-tweet': topic.type === 1,
                }"
                itemprop="articleBody"
              >
                <h1 v-if="topic.title" class="topic-title" itemprop="headline">
                  {{ topic.title }}
                </h1>
                <div
                  ref="topicContent"
                  v-lazy-container="{ selector: 'img' }"
                  class="topic-content-detail"
                  v-html="topic.content"
                ></div>
                <ul
                  v-if="topic.imageList && topic.imageList.length"
                  v-viewer
                  class="topic-image-list"
                >
                  <li v-for="(image, index) in topic.imageList" :key="index">
                    <div class="image-item">
                      <img :src="image.preview" :data-src="image.url" />
                    </div>
                  </li>
                </ul>
              </div>

              <!--节点、标签-->
              <div class="topic-tags">
                <nuxt-link
                  v-if="topic.node"
                  :to="'/topics/node/' + topic.node.nodeId"
                  class="topic-tag"
                  >{{ topic.node.name }}</nuxt-link
                >
                <nuxt-link
                  v-for="tag in topic.tags"
                  :key="tag.tagId"
                  :to="'/topics/tag/' + tag.tagId"
                  class="topic-tag"
                  >#{{ tag.tagName }}</nuxt-link
                >
              </div>

              <!-- 点赞用户列表 -->
              <div
                v-if="likeUsers && likeUsers.length"
                class="topic-like-users"
              >
                <avatar
                  v-for="likeUser in likeUsers"
                  :key="likeUser.id"
                  :user="likeUser"
                  :round="true"
                  size="30"
                />
              </div>

              <!-- 功能按钮 -->
              <div class="topic-actions">
                <a class="action disabled">
                  <i class="action-icon iconfont icon-read" />
                  <span class="content">
                    <span>浏览</span>
                    <span v-if="topic.viewCount > 0">
                      ({{ topic.viewCount }})
                    </span>
                  </span>
                </a>
                <a
                  class="action"
                  :class="{ disabled: liked }"
                  @click="like(topic)"
                >
                  <i
                    class="action-icon iconfont icon-like"
                    :class="{ 'checked-icon': liked }"
                  />
                  <span class="content">
                    <span>点赞</span>
                    <span v-if="topic.likeCount > 0">
                      ({{ topic.likeCount }})
                    </span>
                  </span>
                </a>
                <a class="action" @click="addFavorite(topic.topicId)">
                  <i
                    class="action-icon iconfont"
                    :class="{
                      'icon-has-favorite': favorited,
                      'icon-favorite': !favorited,
                      'checked-icon': favorited,
                    }"
                  />
                  <span class="content">
                    <span>收藏</span>
                  </span>
                </a>
              </div>

              <!-- 评论 -->
              <comment
                :entity-id="topic.topicId"
                :comments-page="commentsPage"
                :comment-count="topic.commentCount"
                :show-ad="false"
                :mode="topic.type === 1 ? 'text' : 'markdown'"
                entity-type="topic"
              />
            </article>
          </div>
        </div>
        <div class="right-container">
          <user-info :user="topic.user" />

          <div class="ad">
            <!-- 展示广告 -->
            <adsbygoogle ad-slot="1742173616" />
          </div>

          <div v-if="topic.toc" ref="toc" class="widget no-bg toc">
            <div class="widget-header">目录</div>
            <div class="widget-content" v-html="topic.toc" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import CommonHelper from '~/common/CommonHelper'

export default {
  async asyncData({ $axios, params, error }) {
    let topic
    try {
      topic = await $axios.get('/api/topic/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '话题不存在',
      })
      return
    }

    const [liked, favorited, commentsPage, likeUsers] = await Promise.all([
      $axios.get('/api/like/liked', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/favorite/favorited', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/topic/recentlikes/' + params.id),
    ])

    return {
      topic,
      commentsPage,
      favorited: favorited.favorited,
      liked: liked.liked,
      likeUsers,
    }
  },
  head() {
    return {
      title: this.$topicSiteTitle(this.topic),
      link: [
        {
          rel: 'stylesheet',
          href: CommonHelper.highlightCss,
        },
      ],
      script: [
        {
          type: 'text/javascript',
          src: CommonHelper.highlightScript,
          callback: () => {
            // 客户端渲染的时候执行这里进行代码高亮
            CommonHelper.initHighlight()
          },
        },
      ],
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
  },
  mounted() {
    // 为了解决服务端渲染时，没有刷新meta中的script，callback没执行，导致代码高亮失败的问题
    // 所以服务端渲染时会调用这里的方法进行代码高亮
    CommonHelper.initHighlight(this)
  },
  methods: {
    async addFavorite(topicId) {
      try {
        if (this.favorited) {
          await this.$axios.get('/api/favorite/delete', {
            params: {
              entityType: 'topic',
              entityId: topicId,
            },
          })
          this.favorited = false
          this.$message.success('已取消收藏')
        } else {
          await this.$axios.get('/api/topic/favorite/' + topicId)
          this.favorited = true
          this.$message.success('收藏成功')
        }
      } catch (e) {
        console.error(e)
        this.$message.error('收藏失败：' + (e.message || e))
      }
    },
    async like(topic) {
      try {
        if (this.liked) {
          return
        }
        await this.$axios.post('/api/topic/like/' + topic.topicId)
        this.liked = true
        topic.likeCount++
        this.likeUsers = this.likeUsers || []
        this.likeUsers.unshift(this.$store.state.user.current)
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.liked = true
          this.$message.error(e.message || e)
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped></style>

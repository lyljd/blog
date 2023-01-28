var CONTENT_KEY = "CACHE_CONTENT"; // 编辑器内容缓存key
var TITLE_KEY = "CACHE_TITLE"; // 标题缓存key
var EXIT_AUTO_SAVE_KEY = "EXIT_AUTO_SAVE"; // 退出时自动保存
var EDIT_EXIT_AUTO_SAVE_KEY = "EDIT_EXIT_AUTO_SAVE"; // 退出时自动保存
var cos = null;
var MdEditor = null;
var headInput = null;
var ArticleItem = {};
ArticleItem.type = 0
ArticleItem.categoryId = 0
var EAS = window.localStorage.getItem(EXIT_AUTO_SAVE_KEY) || "1";
var EEAS = window.localStorage.getItem(EDIT_EXIT_AUTO_SAVE_KEY) || "0";

function setAjaxToken(xhr) {
  xhr.setRequestHeader("Token", localStorage.getItem("token"));
}
function initEditor() {
  // 取默认标题
  headInput.val(ArticleItem.title);
  // 初始化编辑器
  MdEditor = editormd("editormd", {
    width: "99.5%",
    height: window.innerHeight - 65,
    syncScrolling: "single",
    editorTheme: "default",
    path: "/resource/lib/",
    placeholder: "请输入正文",
    appendMarkdown: ArticleItem.markdown,
    codeFold: true,
    saveHTMLToTextarea: true,
    // tocm: true,
    imageUpload: true,
    taskList: true,
    emoji: true,
    imageFormats: ["jpg", "jpeg", "gif", "png", "bmp", "webp"],
    // imageUploadURL: "/api/v1/uploadfile",
    imageUploadCalback: function (files, cb) {
      uploadImage(files[0], cb);
    },
  });
}
function uploadImage(file, cb) {
  let fd = new FormData();
  fd.append("image", file)
  $.ajax({
    url: "/api/v1/image",
    type: "POST",
    data: fd,
    contentType:false,  // ajax上传文件必须；不需使用任何编码
    processData:false,  // ajax上传文件必须；使浏览器不要对数据进行任何处理
    success: function (res) {
      if (res.code === 401) {
        alert(res.error);
        toLogin();
        return;
      }
      if (res.code !== 200) {
        alert(res.error);
        return;
      }
      cb(window.location.protocol+"//"+window.location.host+"/resource/image/"+res.data.fn)
    },
    error: function () {
      alert("上传图片失败，未知错误");
    },
    beforeSend: setAjaxToken,
  });
}

function getArticleItem(id) {
  $(".publish-show").text("更新");
  $(".publish-btn").text("更新");
  $.ajax({
    url: "/api/v1/post/" + id,
    type: "GET",
    contentType: "application/json",
    success: function (res) {
      if (res.code === 401) {
        alert(res.error);
        toLogin();
        return;
      }
      if (res.code !== 200) {
        alert(res.error);
        location.href = "/";
        return;
      }
      ArticleItem = res.data || {};
      initActive();
      initEditor();
    },
    error: function () {
      alert("加载页面失败，未知错误");
      location.href = "/";
    },
    beforeSend: setAjaxToken,
  });
}
function initActive() {
  if (ArticleItem.categoryId && ArticleItem.categoryId>0) {
    $(".update-btn").removeAttr("hidden")
    $(".delete-btn").removeAttr("hidden")
  }
  $(".category li[value=" + ArticleItem.categoryId + "]")
    .addClass("active")
    .siblings()
    .removeClass("active");
  $("#type-box li[value=" + ArticleItem.type + "]")
    .addClass("active")
    .siblings()
    .removeClass("active");
  $(".slug-input").val(ArticleItem.slug);
}
function initCache() {
  if(!localStorage.getItem("token")) {
    alert("请登录")
    location.href = "/"
    return
  }
  headInput = $(".header-input");
  var query = new URLSearchParams(location.search);
  var _id = query.get("id");
  if (_id) {
    if(EEAS !== "1") {
      $(".eas-btn").css({"backgroundColor":"#909399"});
    } else {
      window.addEventListener("beforeunload", saveHandler);
    }
    return getArticleItem(_id);
  } else {
    if(EAS !== "1") {
      $(".eas-btn").css({"backgroundColor":"#909399"});
    } else {
      window.addEventListener("beforeunload", saveHandler);
    }
  }
  // 取本地缓存
  ArticleItem.title = window.localStorage.getItem(TITLE_KEY);
  ArticleItem.markdown = window.localStorage.getItem(CONTENT_KEY) || "";
  // initEditor
  initEditor();
}

function saveHandler() {
  window.localStorage.setItem(TITLE_KEY, headInput.val());
  window.localStorage.setItem(CONTENT_KEY, MdEditor.getMarkdown());
}
function clearHandler() {
  window.localStorage.removeItem(TITLE_KEY);
  window.localStorage.removeItem(CONTENT_KEY);
}

// 提交
function publishHandler() {
  ArticleItem.slug = $(".slug-input").val();
  if (ArticleItem.type === 1 && !ArticleItem.slug)
    return $(".publish-tip").text("请输入自定义链接");
  ArticleItem.title = headInput.val();
  if (!ArticleItem.title) return $(".publish-tip").text("请输入标题");
  ArticleItem.markdown = MdEditor.getMarkdown();
  if (!ArticleItem.markdown) return $(".publish-tip").text("请输入正文");
  ArticleItem.content = MdEditor.getPreviewedHTML();

  $.ajax({
    url: "/api/v1/post",
    type: ArticleItem.pid ? "PUT" : "POST",
    contentType: "application/json",
    data: JSON.stringify(ArticleItem),
    success: function (res) {
      if (res.code === 401) {
        alert(res.error);
        toLogin();
        return;
      }
      if (res.code !== 200) return alert(res.error);
      location.href = "/p/"+res.data.pid;
      clearHandler();
    },
    error: function () {
      alert("提交失败，未知错误")
    },
    beforeSend: setAjaxToken,
  });
}

$(function () {
  // 初始化缓存
  initCache();
  //退出时自动保存
  $(".eas-btn").click(function () {
    if (ArticleItem.pid) {
      if(EEAS === "1") {
        EEAS = "0";
        $(".eas-btn").css({"backgroundColor":"#909399"});
        window.removeEventListener("beforeunload", saveHandler);
      } else {
        EEAS = "1";
        $(".eas-btn").css({"backgroundColor":"#67C23A"});
        window.addEventListener("beforeunload", saveHandler);
      }
      window.localStorage.setItem(EDIT_EXIT_AUTO_SAVE_KEY, EEAS)
    } else {
      if(EAS === "1") {
        EAS = "0";
        $(".eas-btn").css({"backgroundColor":"#909399"});
        window.removeEventListener("beforeunload", saveHandler);
      } else {
        EAS = "1";
        $(".eas-btn").css({"backgroundColor":"#67C23A"});
        window.addEventListener("beforeunload", saveHandler);
      }
      window.localStorage.setItem(EXIT_AUTO_SAVE_KEY, EAS)
    }
  });
  // 返回首页
  var back = $(".home-btn");
  back.click(function () {
    if (ArticleItem.pid) {
      location.href = "/p/" + ArticleItem.pid
    } else {
      let from = new URLSearchParams(location.search).get("from");
      if (from) {
        location.href = from
      } else {
        location.href = "/"
      }
    }
  });
  // 保存
  $(".save-btn").click(function () {
    saveHandler()
    alert("保存成功！")
  });
  var drop = $(".publish-drop");
  // 显示
  $(".publish-show").click(function () {
    drop.show();
  });
  $(".cancel-btn").click(function () {
    drop.hide();
    $(".publish-tip").text("");
  });

  // 新增
  $(".new-btn").click(function () {
    let cn = prompt('请输入分类名')
    if (cn === null) {
      return;
    }
    cn = cn.trim();
    if (cn.length === 0) {
      alert("分类名不能为空");
      return;
    }
    if (cn.length > 10) {
      alert("分类名最长10位");
      return;
    }
    $.ajax({
      url: "/api/v1/category",
      type: "POST",
      data: JSON.stringify({cn: cn}),
      contentType: "application/json",
      success: function (res) {
        if (res.code === 401) {
          alert(res.error);
          toLogin();
          return;
        }
        if (res.code != 200) {
          alert(res.error);
          return;
        }
        $(".new-btn").before("<li value="+res.data.cid+" title="+cn+">"+cn+"</li>");
      },
      error: function () {
        alert("新增失败，未知错误")
      },
      beforeSend: setAjaxToken,
    });
  });
  // 修改
  $(".update-btn").click(function () {
    let cn = prompt('请输入新分类名')
    if (cn === null) {
      return;
    }
    cn = cn.trim();
    if (cn.length === 0) {
      alert("分类名不能为空");
      return;
    }
    if (cn.length > 10) {
      alert("分类名最长10位");
      return;
    }
    $.ajax({
      url: "/api/v1/category/" + ArticleItem.categoryId,
      type: "PUT",
      data: JSON.stringify({cn: cn}),
      contentType: "application/json",
      success: function (res) {
        if (res.code === 401) {
          alert(res.error);
          toLogin();
          return;
        }
        if (res.code != 200) {
          alert(res.error);
          return;
        }
        $(".category li[value=" + ArticleItem.categoryId + "]").html(cn);
      },
      error: function () {
        alert("修改失败，未知错误")
      },
      beforeSend: setAjaxToken,
    });
  });
  //删除
  $(".delete-btn").click(function () {
    var r = confirm("所属此分类的所有文章将会被取消分类，你确定要删除此分类吗？");
    if (!r) return;
    $.ajax({
      url: "/api/v1/category/" + ArticleItem.categoryId,
      type: "DELETE",
      success: function (res) {
        if (res.code === 401) {
          alert(res.error);
          toLogin();
          return;
        }
        if (res.code != 200) {
          alert(res.error);
          return;
        }
        $(".update-btn").attr("hidden",true)
        $(".delete-btn").attr("hidden",true)
        $(".category li[value=" + ArticleItem.categoryId + "]").remove();
        ArticleItem.categoryId = 0
        alert("删除成功")
      },
      error: function () {
        alert("删除失败，未知错误")
      },
      beforeSend: setAjaxToken,
    });
  });

  // 发布逻辑
  $(".publish-btn").click(publishHandler);
  // 选择分类
  $("#category-box").on("click", "li", function (event) {
    var target = $(event.target);
    let seleCid = parseInt(target.attr("value"));
    if (seleCid === ArticleItem.categoryId) {
      ArticleItem.categoryId = 0
      target.removeClass("active");
      $(".update-btn").attr("hidden",true)
      $(".delete-btn").attr("hidden",true)
    } else {
      ArticleItem.categoryId = seleCid
      target.addClass("active").siblings().removeClass("active");
      $(".update-btn").removeAttr("hidden")
      $(".delete-btn").removeAttr("hidden")
    }
    $(".publish-tip").text("");
  });
  // 选择类型
  $("#type-box").on("click", "li", function (event) {
    var target = $(event.target);
    target.addClass("active").siblings().removeClass("active");
    ArticleItem.type = Number(target.attr("value") || 0);
    $(".publish-tip").text("");
  });
});

var TOKEN_KEY = "token";
var USER_KEY = "nickname";
function logout() {
  localStorage.removeItem(USER_KEY);
  localStorage.removeItem(TOKEN_KEY);
  $(".login-action").show();
  $(".login-end").hide();
}
function toLogin() {
  logout();
  location.href = "/login?from=" + window.location.href;
}

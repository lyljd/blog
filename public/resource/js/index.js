var TOKEN_KEY = "token";
var USER_KEY = "nickname";

$(function () {
  // 登录
  loginLogic();
  // 归档
  pagination();
  // toc
  tocInit();
  // 编辑和删除
  initEditLogic();
  // 顶部选中
  headerActive();
  // 归档排序
  initPigeSort();
  // 搜索
  initSearch();
});

function setAjaxToken(xhr) {
  xhr.setRequestHeader("Token", localStorage.getItem(TOKEN_KEY));
}

function headerActive() {
  var nav = $('a[href="' + location.pathname + '"]');
  if (nav.length == 0) return;
  nav.addClass("active");
}
function initEditLogic() {
  var edit = $(".detail-edit");
  if (localStorage.getItem(TOKEN_KEY) && edit.length > 0) {
    edit.show();
    var delEle = $(".detail-delete");
    // 绑定删除事件
    delEle.click(function () {
      deleteDetail(delEle.attr("pid"));
    });
  }
}
// 登录部分逻辑
function loginLogic() {
  if (localStorage.getItem(TOKEN_KEY)) {
    $(".login-action").hide();
    $(".login-end").show();
    $(".login-username").text(localStorage.getItem(USER_KEY));
  }
  // 登录
  $(".login-submint").click(function () {
    var tipEle = $(".login-tip");
    var name = $(".login-name").val();
    var passwd = $(".login-passwd").val();
    if (!name) return tipEle.show().text("请输入用户名");
    if (!passwd) return tipEle.show().text("请输入密码");
    tipEle.show().text("请稍等...");
    $.ajax({
      url: "/api/v1/login",
      data: JSON.stringify({ username: name, password: passwd }),
      contentType: "application/json",
      type: "POST",
      success: function (res) {
          if (res.code !== 200) {
            return tipEle.show().text(res.error);
          }
          localStorage.setItem(TOKEN_KEY, res.data.token);
          localStorage.setItem(USER_KEY, res.data.nickname);
          let index = window.location.href.indexOf("from=")
          location.href = window.location.href.slice(index+5);
      },
      error: function (err) {
        console.log("err", err);
        tipEle.show().text("登录失败，未知错误");
      },
    });
  });
  // 退出登录
  $(".login-out").click(function () {
    logout();
    location.reload();
  });
}

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

// 翻页逻辑
function pagination() {
  var query = new URLSearchParams(location.search);
  var page = query.get("page") || 1;
  $(".pagination-next").click(function () {
    page++;
    location.search = "?page=" + page;
  });
  $(".pagination-prev").click(function () {
    page--;
    location.search = "?page=" + page;
  });
  // $(".pagination-btn").click(function (event) {
  //   var val = $(event.target).attr("value");
  //   if (val == 1) return (location.href = "/");
  //   location.search = "?page=" + val;
  // });
}
function deleteDetail(id) {
  var r = confirm("你确定要永久地删除此文章吗？");
  if (!r) return;
  $.ajax({
    url: "/api/v1/post/" + id,
    type: "DELETE",
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
      alert("删除成功")
      location.href = "/";
    },
    error: function () {
      alert("删除失败，未知错误")
    },
    beforeSend: setAjaxToken,
  });
}
function tocInit() {
  var tocBox = $("#toc-box");
  if (tocBox.length == 0) return;
  imageZoom();
  var mdTocList = $(".markdown-toc-list");
  // 如果有TOC
  if (mdTocList.length > 0 && mdTocList.children().length > 0) {
    tocBox.append(mdTocList);
    tocScrollTo(tocBox);
  } else {
    $(".detail-left").css("width", "100%");
    $(".detail-right").hide();
  }
}
function imageZoom() {
  var zoom = $(".zoom-prev");
  $(".detail-content").on("click", "img", function (event) {
    $(".zoom-container").css(
      "background-image",
      "url(" + $(event.target).attr("src") + ")"
    );
    zoom.show();
  });
  zoom.click(function () {
    zoom.hide();
  });
}
function tocScrollTo(tocBox) {
  // 组织默认事件
  var all = document.querySelectorAll("#toc-box a");
  for (var i = 0, len = all.length; i < len; i++) {
    all[i].href = "javascript:void(0)";
  }
  var prvEle = null;
  tocBox.on("click", "a", function (event) {
    event.stopPropagation();
    ele = $(event.target);
    ele.addClass("active");
    if (prvEle) prvEle.removeClass("active");
    prvEle = ele;
    var _href = $(event.target).text();
    var top = $("a[name='" + _href + "']").offset().top;
    window.scrollTo(0, top - 80);
  });
}

function initPigeSort() {
  var box = $(".pige-content");
  if (box.length == 0) return;
  var children = box.children();
  // 翻转排序
  for (var i = children.length; i >= 0; i--) {
    box.append(children[i]);
  }
}

function initSearch() {
  var timer = null;
  var timer2 = null;
  var searchList = [];
  var drop = $(".search-drop");
  var input = $("#search-input");
  input.on("input", function (event) {
    clearTimeout(timer);
    timer = setTimeout(function () {
      searchHandler(event.target.value);
    }, 300);
  });
  input.on("blur", function () {
    timer2 = setTimeout(function () {
      drop.hide();
    }, 100);
  });
  drop.on("click", function () {
    clearTimeout(timer2);
  });
  input.on("focus", function (event) {
    if (searchList.length > 0) {
      drop.show().html(searchList.join(""));
    } else {
      searchHandler(event.target.value);
    }
  });
  function searchHandler(val) {
    if (!val) {
      drop.hide();
      searchList = [];
      return;
    };
    $.ajax({
      url: "/api/v1/post/search?k=" + val,
      contentType: "application/json",
      success: function (res) {
        if (res.code !== 200){
          searchList = ["<span>无搜索结果<span/>"];
          drop.show().html(searchList.join(""));
          return
        }
        var data = res.data || [];
        searchList = [];
        if (data.length === 0) {
          searchList = ["<span>无搜索结果<span/>"];
          drop.show().html(searchList.join(""));
          return
        }
        for (var i = 0, len = data.length; i < len; i++) {
          var item = data[i];
          searchList.push("<a href='/p/" + item.pid + "'>" + item.title + "<a/>");
          drop.show().html(searchList.join(""));
        }
      },
      error: function () {
        searchList = ["<span>搜索失败<span/>"];
        drop.show().html(searchList.join(""));
      },
    });
  }
}

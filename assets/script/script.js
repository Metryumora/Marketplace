/**
 * Created by Metr_yumora on 27.04.2017.
 */
function hideText() {
    $('p.description').html(function (i, t) {
        return t.replace("/e/g", 'Z');
    })
}

//return t.replace("^$", '<span class="hidden">Hello</span>');
#define A foo
#define B bar
#define J(X,Y) X##Y
x = "A##B";
y = J(A,B);

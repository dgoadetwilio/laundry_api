alias kc="kubectl"
alias kctoken='kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '\''{print $1}'\'')'
alias cc='kubectl config current-context'
alias uc='kubectl config use-context'
alias deletepodsinerrorstate='kubectl get pods --field-selector=status.phase=Error | awk "{print \$1}" | xargs -I {} kubectl delete pods/{}'
alias kcimages='kubectl get pods -o jsonpath="{.items[*].spec.containers[*].image}" |tr -s "[[:space:]]" "\n" |sort |uniq -c'
